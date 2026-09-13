package roaringv3

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRankAndSelect(t *testing.T) {
	t.Parallel()

	type args struct {
		values      []uint32
		runOptimize bool
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "empty", args: args{values: nil}},
		{name: "array container", args: args{values: []uint32{3, 9, 100, 65535}}},
		{name: "bitmap container", args: args{values: spread(0, 3, 9000)}},
		{name: "run container", args: args{values: sequence(0, 20000), runOptimize: true}},
		{name: "several chunks", args: args{values: append(sequence(0, 9000), sequence(1<<20, 50)...)}},
		{name: "runs across several chunks", args: args{values: append(sequence(0, 20000), sequence(1<<20, 20000)...), runOptimize: true}},
		{name: "top of the key space", args: args{values: []uint32{^uint32(0) - 1, ^uint32(0)}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := New(tt.args.values...)
			if tt.args.runOptimize {
				b.RunOptimize()
			}
			want := b.ToArray()

			for i, v := range want {
				require.Equal(t, i+1, b.Rank(v), "Rank(%d)", v)

				got, ok := b.Select(i)
				require.True(t, ok, "Select(%d)", i)
				require.Equal(t, v, got, "Select(%d)", i)
			}

			_, ok := b.Select(len(want))
			require.False(t, ok, "Select past the end")
			_, ok = b.Select(-1)
			require.False(t, ok, "Select of a negative index")

			gotMin, okMin := b.Minimum()
			gotMax, okMax := b.Maximum()
			require.Equal(t, len(want) > 0, okMin)
			require.Equal(t, len(want) > 0, okMax)
			if len(want) > 0 {
				require.Equal(t, want[0], gotMin)
				require.Equal(t, want[len(want)-1], gotMax)
			}
		})
	}
}

func TestRankOnValuesOutsideTheSet(t *testing.T) {
	t.Parallel()

	type args struct {
		query uint32
	}

	b := New(10, 20, 1<<16, 1<<16+5)

	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "below everything", args: args{query: 0}, want: 0},
		{name: "between two values", args: args{query: 15}, want: 1},
		{name: "at the end of a chunk", args: args{query: 65535}, want: 2},
		{name: "inside a later chunk", args: args{query: 1<<16 + 2}, want: 3},
		{name: "above everything", args: args{query: ^uint32(0)}, want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, b.Rank(tt.args.query))
		})
	}
}

func FuzzRankSelectAgainstReference(f *testing.F) {
	f.Add(uint64(1), uint32(1000), false)
	f.Add(uint64(5), uint32(70000), true)

	f.Fuzz(func(t *testing.T, seed uint64, span uint32, runOptimize bool) {
		rng := rand.New(rand.NewPCG(seed, seed+1))
		values := randomRuns(rng, int(span%300), span|1)

		b := New(values...)
		if runOptimize {
			b.RunOptimize()
		}
		ref := sortedUnique(values)

		for i, v := range ref {
			got, ok := b.Select(i)
			require.True(t, ok)
			require.Equal(t, v, got)
			require.Equal(t, i+1, b.Rank(v))
		}

		// Rank of an arbitrary probe must match a linear count.
		probe := rng.Uint32N(span | 1)
		wantRank, _ := slices.BinarySearch(ref, probe)
		if idx, found := slices.BinarySearch(ref, probe); found {
			wantRank = idx + 1
		}
		require.Equal(t, wantRank, b.Rank(probe))
	})
}
