//go:build goexperiment.simd && amd64

package simdops

import (
	"math/bits"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitwiseTo(t *testing.T) {
	t.Parallel()

	type args struct {
		a []uint64
		b []uint64
	}

	tests := []struct {
		name string
		args args
		want map[string][]uint64
	}{
		{
			name: "empty",
			args: args{a: []uint64{}, b: []uint64{}},
			want: map[string][]uint64{"and": {}, "or": {}, "andnot": {}, "xor": {}},
		},
		{
			name: "shorter than one vector",
			args: args{a: []uint64{0b1100, 0b1010, 0b0001}, b: []uint64{0b1010, 0b0110, 0b0001}},
			want: map[string][]uint64{
				"and":    {0b1000, 0b0010, 0b0001},
				"or":     {0b1110, 0b1110, 0b0001},
				"andnot": {0b0100, 0b1000, 0b0000},
				"xor":    {0b0110, 0b1100, 0b0000},
			},
		},
		{
			name: "exactly one vector",
			args: args{
				a: []uint64{1, 2, 3, 4, 5, 6, 7, 8},
				b: []uint64{8, 7, 6, 5, 4, 3, 2, 1},
			},
			want: map[string][]uint64{
				"and":    {0, 2, 2, 4, 4, 2, 2, 0},
				"or":     {9, 7, 7, 5, 5, 7, 7, 9},
				"andnot": {1, 0, 1, 0, 1, 4, 5, 8},
				"xor":    {9, 5, 5, 1, 1, 5, 5, 9},
			},
		},
	}

	ops := map[string]func(dst, a, b []uint64) int{
		"and":    AndTo,
		"or":     OrTo,
		"andnot": AndNotTo,
		"xor":    XorTo,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for name, op := range ops {
				dst := make([]uint64, len(tt.args.a))
				got := op(dst, tt.args.a, tt.args.b)

				require.Equal(t, tt.want[name], dst, name)
				require.Equal(t, popcountRef(tt.want[name]), got, "%s cardinality", name)
			}
		})
	}
}

func TestBitwiseToMatchesReferenceOnFullContainer(t *testing.T) {
	t.Parallel()

	type args struct {
		density float64
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "sparse", args: args{density: 0.01}},
		{name: "half full", args: args{density: 0.5}},
		{name: "dense", args: args{density: 0.99}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rng := rand.New(rand.NewPCG(1, 2))
			a := randomBitmap(rng, tt.args.density)
			b := randomBitmap(rng, tt.args.density)

			wantAnd := make([]uint64, BitmapWords)
			wantOr := make([]uint64, BitmapWords)
			wantAndNot := make([]uint64, BitmapWords)
			wantXor := make([]uint64, BitmapWords)
			for i := range a {
				wantAnd[i] = a[i] & b[i]
				wantOr[i] = a[i] | b[i]
				wantAndNot[i] = a[i] &^ b[i]
				wantXor[i] = a[i] ^ b[i]
			}

			dst := make([]uint64, BitmapWords)

			require.Equal(t, popcountRef(wantAnd), AndTo(dst, a, b))
			require.Equal(t, wantAnd, dst)
			require.Equal(t, popcountRef(wantOr), OrTo(dst, a, b))
			require.Equal(t, wantOr, dst)
			require.Equal(t, popcountRef(wantAndNot), AndNotTo(dst, a, b))
			require.Equal(t, wantAndNot, dst)
			require.Equal(t, popcountRef(wantXor), XorTo(dst, a, b))
			require.Equal(t, wantXor, dst)

			require.Equal(t, popcountRef(wantAnd), AndCardinality(a, b))
			require.Equal(t, popcountRef(wantOr), OrCardinality(a, b))
			require.Equal(t, popcountRef(wantAndNot), AndNotCardinality(a, b))
			require.Equal(t, popcountRef(wantXor), XorCardinality(a, b))
			require.Equal(t, popcountRef(wantAnd) > 0, Intersects(a, b))
		})
	}
}

func TestAndToAliasesInput(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewPCG(7, 7))
	a := randomBitmap(rng, 0.3)
	b := randomBitmap(rng, 0.3)

	want := make([]uint64, BitmapWords)
	for i := range a {
		want[i] = a[i] & b[i]
	}

	got := AndTo(a, a, b)

	require.Equal(t, want, a)
	require.Equal(t, popcountRef(want), got)
}

func TestPopcount(t *testing.T) {
	t.Parallel()

	type args struct {
		words []uint64
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "empty", args: args{words: nil}, want: 0},
		{name: "single word", args: args{words: []uint64{0b1011}}, want: 3},
		{name: "tail shorter than a vector", args: args{words: []uint64{^uint64(0), 1, 3}}, want: 67},
		{name: "all ones", args: args{words: onesWords(BitmapWords)}, want: BitmapWords * 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, Popcount(tt.args.words))
		})
	}
}

func TestIntersects(t *testing.T) {
	t.Parallel()

	type args struct {
		a []uint64
		b []uint64
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "disjoint", args: args{a: []uint64{0b0101}, b: []uint64{0b1010}}, want: false},
		{name: "overlapping", args: args{a: []uint64{0b0110}, b: []uint64{0b0010}}, want: true},
		{name: "hit past the first vector", args: args{a: bitAt(600), b: bitAt(600)}, want: true},
		{name: "miss across a full container", args: args{a: bitAt(600), b: bitAt(601)}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, Intersects(tt.args.a, tt.args.b))
		})
	}
}

func TestBitmapToArray(t *testing.T) {
	t.Parallel()

	type args struct {
		values []uint16
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "empty", args: args{values: []uint16{}}},
		{name: "first and last bit", args: args{values: []uint16{0, 65535}}},
		{name: "within one word", args: args{values: []uint16{1, 5, 63}}},
		{name: "crossing the 32-bit halves", args: args{values: []uint16{31, 32, 33}}},
		{name: "spread across the container", args: args{values: []uint16{0, 64, 1000, 4095, 40000, 65535}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bm := make([]uint64, BitmapWords)
			for _, v := range tt.args.values {
				bm[v/64] |= 1 << (v % 64)
			}

			dst := make([]uint16, len(tt.args.values))
			n := BitmapToArray(dst, bm)

			require.Equal(t, len(tt.args.values), n)
			require.Equal(t, tt.args.values, dst[:n])
		})
	}
}

func TestBitmapToArrayMatchesReference(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewPCG(11, 13))
	bm := randomBitmap(rng, 0.2)

	want := make([]uint16, 0, Popcount(bm))
	for i := 0; i < BitmapWords*64; i++ {
		if bm[i/64]&(1<<(i%64)) != 0 {
			want = append(want, uint16(i))
		}
	}

	dst := make([]uint16, len(want))
	n := BitmapToArray(dst, bm)

	require.Equal(t, len(want), n)
	require.Equal(t, want, dst[:n])
}

func popcountRef(words []uint64) int {
	n := 0
	for _, w := range words {
		n += bits.OnesCount64(w)
	}
	return n
}

func randomBitmap(rng *rand.Rand, density float64) []uint64 {
	bm := make([]uint64, BitmapWords)
	for i := 0; i < BitmapWords*64; i++ {
		if rng.Float64() < density {
			bm[i/64] |= 1 << (i % 64)
		}
	}
	return bm
}

func onesWords(n int) []uint64 {
	w := make([]uint64, n)
	for i := range w {
		w[i] = ^uint64(0)
	}
	return w
}

func bitAt(pos int) []uint64 {
	bm := make([]uint64, BitmapWords)
	bm[pos/64] |= 1 << (pos % 64)
	return bm
}
