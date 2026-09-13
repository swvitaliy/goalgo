package simdops

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntersectArrays(t *testing.T) {
	t.Parallel()

	type args struct {
		a []uint16
		b []uint16
	}

	tests := []struct {
		name string
		args args
		want []uint16
	}{
		{name: "both empty", args: args{a: []uint16{}, b: []uint16{}}, want: []uint16{}},
		{name: "one empty", args: args{a: []uint16{1, 2, 3}, b: []uint16{}}, want: []uint16{}},
		{name: "disjoint", args: args{a: []uint16{1, 3, 5}, b: []uint16{2, 4, 6}}, want: []uint16{}},
		{name: "identical", args: args{a: []uint16{1, 2, 3}, b: []uint16{1, 2, 3}}, want: []uint16{1, 2, 3}},
		{name: "partial overlap", args: args{a: []uint16{1, 4, 7, 9}, b: []uint16{2, 4, 9}}, want: []uint16{4, 9}},
		{
			name: "block-wise merge across several vectors",
			args: args{a: rangeVals(0, 200, 1), b: rangeVals(0, 200, 2)},
			want: rangeVals(0, 200, 2),
		},
		{
			name: "merge where blocks advance unevenly",
			args: args{a: rangeVals(0, 300, 1), b: rangeVals(150, 450, 1)},
			want: rangeVals(150, 300, 1),
		},
		{
			name: "galloping when one side is far larger",
			args: args{a: []uint16{0, 5000, 20000, 65535}, b: rangeVals(0, 30000, 1)},
			want: []uint16{0, 5000, 20000},
		},
		{
			name: "galloping with the small array first",
			args: args{a: rangeVals(0, 30000, 1), b: []uint16{3, 17, 29999, 30000}},
			want: []uint16{3, 17, 29999},
		},
		{
			name: "galloping hit on the very first probe",
			args: args{a: []uint16{0}, b: rangeVals(0, 30000, 1)},
			want: []uint16{0},
		},
		{
			name: "galloping miss below the whole range",
			args: args{a: []uint16{7}, b: rangeVals(100, 30100, 1)},
			want: []uint16{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dst := make([]uint16, min(len(tt.args.a), len(tt.args.b)))
			n := IntersectArrays(dst, tt.args.a, tt.args.b)

			require.Equal(t, tt.want, dst[:n])
		})
	}
}

func TestIntersectArraysMatchesReference(t *testing.T) {
	t.Parallel()

	type args struct {
		sizeA int
		sizeB int
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "equal sizes take the merge path", args: args{sizeA: 2000, sizeB: 2000}},
		{name: "mild imbalance stays on the merge path", args: args{sizeA: 500, sizeB: 2000}},
		{name: "heavy imbalance takes the gallop path", args: args{sizeA: 40, sizeB: 20000}},
		{name: "sizes around a single vector", args: args{sizeA: 31, sizeB: 33}},
		{name: "one value against many", args: args{sizeA: 1, sizeB: 5000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rng := rand.New(rand.NewPCG(uint64(tt.args.sizeA), uint64(tt.args.sizeB)))
			a := randomArray(rng, tt.args.sizeA)
			b := randomArray(rng, tt.args.sizeB)

			dst := make([]uint16, min(len(a), len(b)))
			n := IntersectArrays(dst, a, b)

			require.Equal(t, intersectRef(a, b), dst[:n])
		})
	}
}

func TestArrayMerges(t *testing.T) {
	t.Parallel()

	type args struct {
		a []uint16
		b []uint16
	}

	tests := []struct {
		name          string
		args          args
		wantUnion     []uint16
		wantDifferene []uint16
		wantXor       []uint16
	}{
		{
			name:          "both empty",
			args:          args{a: []uint16{}, b: []uint16{}},
			wantUnion:     []uint16{},
			wantDifferene: []uint16{},
			wantXor:       []uint16{},
		},
		{
			name:          "disjoint",
			args:          args{a: []uint16{1, 3}, b: []uint16{2, 4}},
			wantUnion:     []uint16{1, 2, 3, 4},
			wantDifferene: []uint16{1, 3},
			wantXor:       []uint16{1, 2, 3, 4},
		},
		{
			name:          "partial overlap",
			args:          args{a: []uint16{1, 2, 5, 9}, b: []uint16{2, 5, 7}},
			wantUnion:     []uint16{1, 2, 5, 7, 9},
			wantDifferene: []uint16{1, 9},
			wantXor:       []uint16{1, 7, 9},
		},
		{
			name:          "b exhausted first",
			args:          args{a: []uint16{1, 2, 3, 4}, b: []uint16{1}},
			wantUnion:     []uint16{1, 2, 3, 4},
			wantDifferene: []uint16{2, 3, 4},
			wantXor:       []uint16{2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dst := make([]uint16, len(tt.args.a)+len(tt.args.b))

			n := UnionArrays(dst, tt.args.a, tt.args.b)
			require.Equal(t, tt.wantUnion, dst[:n])

			n = DifferenceArrays(dst, tt.args.a, tt.args.b)
			require.Equal(t, tt.wantDifferene, dst[:n])

			n = XorArrays(dst, tt.args.a, tt.args.b)
			require.Equal(t, tt.wantXor, dst[:n])
		})
	}
}

func FuzzIntersectArrays(f *testing.F) {
	f.Add([]byte{1, 0, 2, 0}, []byte{2, 0, 3, 0})
	f.Add([]byte{}, []byte{1, 0})

	f.Fuzz(func(t *testing.T, rawA, rawB []byte) {
		a := sortedSet(rawA)
		b := sortedSet(rawB)

		dst := make([]uint16, min(len(a), len(b)))
		n := IntersectArrays(dst, a, b)

		require.Equal(t, intersectRef(a, b), dst[:n])
	})
}

func intersectRef(a, b []uint16) []uint16 {
	set := make(map[uint16]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}

	out := make([]uint16, 0, min(len(a), len(b)))
	for _, v := range a {
		if _, ok := set[v]; ok {
			out = append(out, v)
		}
	}
	return out
}

func rangeVals(lo, hi, step int) []uint16 {
	out := make([]uint16, 0, (hi-lo)/step)
	for v := lo; v < hi; v += step {
		out = append(out, uint16(v))
	}
	return out
}

func randomArray(rng *rand.Rand, n int) []uint16 {
	set := make(map[uint16]struct{}, n)
	for len(set) < n {
		set[uint16(rng.UintN(1<<16))] = struct{}{}
	}

	out := make([]uint16, 0, n)
	for v := range set {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

// sortedSet turns fuzz bytes into a sorted, duplicate-free uint16 array.
func sortedSet(raw []byte) []uint16 {
	set := make(map[uint16]struct{}, len(raw)/2)
	for i := 0; i+1 < len(raw); i += 2 {
		set[uint16(raw[i])|uint16(raw[i+1])<<8] = struct{}{}
	}

	out := make([]uint16, 0, len(set))
	for v := range set {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}
