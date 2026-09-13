//go:build goexperiment.simd && amd64

package simdops

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRangeOps(t *testing.T) {
	t.Parallel()

	type args struct {
		start int
		end   int
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "empty range", args: args{start: 100, end: 100}},
		{name: "single bit", args: args{start: 7, end: 8}},
		{name: "within one word", args: args{start: 3, end: 40}},
		{name: "exactly one word", args: args{start: 64, end: 128}},
		{name: "crossing a word boundary", args: args{start: 60, end: 70}},
		{name: "spanning many words", args: args{start: 5, end: 1000}},
		{name: "spanning several vectors", args: args{start: 1, end: 40000}},
		{name: "up to the last bit", args: args{start: 65000, end: BitmapWords * 64}},
		{name: "the whole container", args: args{start: 0, end: BitmapWords * 64}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			want := make([]uint64, BitmapWords)
			for i := tt.args.start; i < tt.args.end; i++ {
				want[i/64] |= 1 << (i % 64)
			}

			set := make([]uint64, BitmapWords)
			SetRange(set, tt.args.start, tt.args.end)
			require.Equal(t, want, set, "SetRange")

			cleared := onesWords(BitmapWords)
			ClearRange(cleared, tt.args.start, tt.args.end)
			for i := range cleared {
				require.Equal(t, ^want[i], cleared[i], "ClearRange word %d", i)
			}

			flipped := make([]uint64, BitmapWords)
			FlipRange(flipped, tt.args.start, tt.args.end)
			require.Equal(t, want, flipped, "FlipRange from zero")
			FlipRange(flipped, tt.args.start, tt.args.end)
			require.Equal(t, make([]uint64, BitmapWords), flipped, "FlipRange is its own inverse")

			src := onesWords(BitmapWords)
			dst := make([]uint64, BitmapWords)
			CopyRange(dst, src, tt.args.start, tt.args.end)
			require.Equal(t, want, dst, "CopyRange")
		})
	}
}

func TestCopyRangeLeavesSurroundingBitsAlone(t *testing.T) {
	t.Parallel()

	src := make([]uint64, BitmapWords)
	SetRange(src, 100, 200)

	dst := onesWords(BitmapWords)
	CopyRange(dst, src, 190, 210)

	// Inside the copied range dst now mirrors src, so 200..209 went from set to
	// clear; everything outside stays as it was.
	want := onesWords(BitmapWords)
	ClearRange(want, 200, 210)
	require.Equal(t, want, dst)
}
