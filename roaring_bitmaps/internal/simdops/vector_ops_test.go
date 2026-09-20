//go:build goexperiment.simd && amd64

package simdops

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVectorOpsMatchesPackageFunctions(t *testing.T) {
	t.Parallel()

	require.True(t, Vectorized)

	type args struct {
		a     []uint64
		b     []uint64
		arrA  []uint16
		arrB  []uint16
		start int
		end   int
	}

	rng := rand.New(rand.NewPCG(7, 11))
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sparse bitmaps and small arrays",
			args: args{
				a: randomBitmap(rng, 0.01), b: randomBitmap(rng, 0.02),
				arrA: randomArray(rng, 50), arrB: randomArray(rng, 70),
				start: 100, end: 60000,
			},
		},
		{
			name: "dense bitmaps and large arrays",
			args: args{
				a: randomBitmap(rng, 0.7), b: randomBitmap(rng, 0.5),
				arrA: randomArray(rng, 4000), arrB: randomArray(rng, 3000),
				start: 63, end: 65, // range inside one word
			},
		},
		{
			name: "full and empty bitmaps",
			args: args{
				a: onesWords(BitmapWords), b: make([]uint64, BitmapWords),
				arrA: rangeVals(0, 65536, 3), arrB: rangeVals(1, 65536, 5),
				start: 0, end: BitmapWords * 64,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var ops VectorOps
			checkWordOps(t, ops, tt.args.a, tt.args.b, tt.args.start, tt.args.end)
			checkArrayOps(t, ops, tt.args.arrA, tt.args.arrB)
			checkRangeOps(t, ops, tt.args.a, tt.args.b, tt.args.start, tt.args.end)
		})
	}
}

func checkWordOps(t *testing.T, ops VectorOps, a, b []uint64, start, end int) {
	t.Helper()

	want := make([]uint64, BitmapWords)
	got := make([]uint64, BitmapWords)
	require.Equal(t, AndTo(want, a, b), ops.AndTo(got, a, b))
	require.Equal(t, want, got)
	require.Equal(t, OrTo(want, a, b), ops.OrTo(got, a, b))
	require.Equal(t, want, got)
	require.Equal(t, AndNotTo(want, a, b), ops.AndNotTo(got, a, b))
	require.Equal(t, want, got)
	require.Equal(t, XorTo(want, a, b), ops.XorTo(got, a, b))
	require.Equal(t, want, got)

	require.Equal(t, AndCardinality(a, b), ops.AndCardinality(a, b))
	require.Equal(t, Popcount(a), ops.Popcount(a))
	require.Equal(t, PopcountPrefix(a, end), ops.PopcountPrefix(a, end))
	require.Equal(t, SelectBit(a, Popcount(a)/2), ops.SelectBit(a, Popcount(a)/2))
	require.Equal(t, Intersects(a, b), ops.Intersects(a, b))
	require.Equal(t, HasBitInRange(b, start, end), ops.HasBitInRange(b, start, end))

	wantArr := make([]uint16, BitmapWords*64)
	gotArr := make([]uint16, BitmapWords*64)
	require.Equal(t, BitmapToArray(wantArr, a), ops.BitmapToArray(gotArr, a))
	require.Equal(t, wantArr, gotArr)
}

func checkArrayOps(t *testing.T, ops VectorOps, a, b []uint16) {
	t.Helper()

	want := make([]uint16, len(a)+len(b))
	got := make([]uint16, len(a)+len(b))
	require.Equal(t, IntersectArrays(want, a, b), ops.IntersectArrays(got, a, b))
	require.Equal(t, want, got)
	require.Equal(t, UnionArrays(want, a, b), ops.UnionArrays(got, a, b))
	require.Equal(t, want, got)
	require.Equal(t, DifferenceArrays(want, a, b), ops.DifferenceArrays(got, a, b))
	require.Equal(t, want, got)
	require.Equal(t, XorArrays(want, a, b), ops.XorArrays(got, a, b))
	require.Equal(t, want, got)
}

func checkRangeOps(t *testing.T, ops VectorOps, src, bm []uint64, start, end int) {
	t.Helper()

	want := append([]uint64(nil), bm...)
	got := append([]uint64(nil), bm...)
	SetRange(want, start, end)
	ops.SetRange(got, start, end)
	require.Equal(t, want, got)
	ClearRange(want, start, end)
	ops.ClearRange(got, start, end)
	require.Equal(t, want, got)
	FlipRange(want, start, end)
	ops.FlipRange(got, start, end)
	require.Equal(t, want, got)
	CopyRange(want, src, start, end)
	ops.CopyRange(got, src, start, end)
	require.Equal(t, want, got)
}
