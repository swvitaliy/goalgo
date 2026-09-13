//go:build goexperiment.simd && amd64

package roaringv1

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddContainsRemove(t *testing.T) {
	t.Parallel()

	type args struct {
		values []uint32
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "empty", args: args{values: nil}},
		{name: "single value", args: args{values: []uint32{42}}},
		{name: "boundaries of the key space", args: args{values: []uint32{0, 65535, 65536, 1 << 31, ^uint32(0)}}},
		{name: "sparse across many chunks", args: args{values: spread(0, 1<<20, 4096)}},
		{name: "dense enough to become a bitmap container", args: args{values: sequence(0, 5000)}},
		{name: "dense in a high chunk", args: args{values: sequence(1<<20, 5000)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := New(tt.args.values...)

			require.Equal(t, len(tt.args.values), b.Cardinality())
			require.Equal(t, len(tt.args.values) == 0, b.IsEmpty())
			require.Equal(t, sortedUnique(tt.args.values), b.ToArray())

			for _, v := range tt.args.values {
				require.True(t, b.Contains(v), "missing %d", v)
				require.False(t, b.Add(v), "re-adding %d should not change the set", v)
			}

			for _, v := range tt.args.values {
				require.True(t, b.Remove(v), "removing %d should change the set", v)
				require.False(t, b.Contains(v), "still present after removal: %d", v)
				require.False(t, b.Remove(v), "double removal of %d should not change the set", v)
			}

			require.True(t, b.IsEmpty())
			require.Equal(t, []uint32{}, b.ToArray())
		})
	}
}

func TestContainerConversionRoundTrip(t *testing.T) {
	t.Parallel()

	b := New(sequence(0, arrayMax)...)
	require.Equal(t, kindArray, b.conts[0].kind)

	require.True(t, b.Add(arrayMax))
	require.Equal(t, kindBitmap, b.conts[0].kind)
	require.Equal(t, arrayMax+1, b.Cardinality())

	require.True(t, b.Remove(arrayMax))
	require.Equal(t, kindArray, b.conts[0].kind)
	require.Equal(t, sequence(0, arrayMax), b.ToArray())
}

func TestSetOperations(t *testing.T) {
	t.Parallel()

	type args struct {
		a []uint32
		b []uint32
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "both empty", args: args{}},
		{name: "left empty", args: args{b: sequence(0, 100)}},
		{name: "right empty", args: args{a: sequence(0, 100)}},
		{name: "array against array, overlapping", args: args{a: sequence(0, 100), b: sequence(50, 100)}},
		{name: "array against array, disjoint", args: args{a: sequence(0, 100), b: sequence(1000, 100)}},
		{name: "array against bitmap", args: args{a: sequence(0, 100), b: sequence(0, 9000)}},
		{name: "bitmap against array", args: args{a: sequence(0, 9000), b: sequence(8000, 100)}},
		{name: "bitmap against bitmap", args: args{a: sequence(0, 9000), b: sequence(4000, 9000)}},
		{name: "disjoint chunks", args: args{a: sequence(0, 9000), b: sequence(1<<20, 9000)}},
		{name: "many sparse chunks", args: args{a: spread(0, 1<<22, 3000), b: spread(1, 1<<22, 3000)}},
		{name: "one chunk deep in the key space", args: args{a: sequence(^uint32(0)-5000, 5000), b: sequence(^uint32(0)-2500, 2500)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a, b := New(tt.args.a...), New(tt.args.b...)
			refA, refB := sortedUnique(tt.args.a), sortedUnique(tt.args.b)

			require.Equal(t, intersectRef(refA, refB), a.And(b).ToArray(), "And")
			require.Equal(t, unionRef(refA, refB), a.Or(b).ToArray(), "Or")
			require.Equal(t, differenceRef(refA, refB), a.AndNot(b).ToArray(), "AndNot")
			require.Equal(t, xorRef(refA, refB), a.Xor(b).ToArray(), "Xor")

			require.Equal(t, len(intersectRef(refA, refB)) > 0, a.Intersects(b), "Intersects")
			require.Equal(t, len(intersectRef(refA, refB)), a.AndCardinality(b), "AndCardinality")

			// Operands must be left untouched.
			require.Equal(t, refA, a.ToArray())
			require.Equal(t, refB, b.ToArray())
		})
	}
}

func TestCloneIsIndependent(t *testing.T) {
	t.Parallel()

	b := New(sequence(0, 9000)...)
	clone := b.Clone()

	require.True(t, clone.Remove(0))
	require.True(t, clone.Add(1<<20))

	require.True(t, b.Contains(0))
	require.False(t, b.Contains(1<<20))
	require.Equal(t, 9000, b.Cardinality())
}

func TestString(t *testing.T) {
	t.Parallel()

	require.Equal(t, "{}", New().String())
	require.Equal(t, "{0,7,65536}", New(65536, 0, 7).String())
}

func FuzzBitmapAgainstReference(f *testing.F) {
	f.Add(uint64(1), uint32(64), uint32(200))
	f.Add(uint64(9), uint32(20000), uint32(70000))

	f.Fuzz(func(t *testing.T, seed uint64, spanA, spanB uint32) {
		rng := rand.New(rand.NewPCG(seed, seed+1))
		valsA := randomValues(rng, int(spanA%20000), spanA|1)
		valsB := randomValues(rng, int(spanB%20000), spanB|1)

		a, b := New(valsA...), New(valsB...)
		refA, refB := sortedUnique(valsA), sortedUnique(valsB)

		require.Equal(t, refA, a.ToArray())
		require.Equal(t, intersectRef(refA, refB), a.And(b).ToArray())
		require.Equal(t, unionRef(refA, refB), a.Or(b).ToArray())
		require.Equal(t, differenceRef(refA, refB), a.AndNot(b).ToArray())
		require.Equal(t, xorRef(refA, refB), a.Xor(b).ToArray())
		require.Equal(t, len(intersectRef(refA, refB)), a.AndCardinality(b))
	})
}

func sequence(start uint32, n int) []uint32 {
	out := make([]uint32, n)
	for i := range out {
		out[i] = start + uint32(i)
	}
	return out
}

// spread returns n values starting at start, stepping far enough apart to land
// in distinct chunks.
func spread(start, step uint32, n int) []uint32 {
	out := make([]uint32, n)
	for i := range out {
		out[i] = start + uint32(i)*step
	}
	return out
}

func randomValues(rng *rand.Rand, n int, span uint32) []uint32 {
	out := make([]uint32, n)
	for i := range out {
		out[i] = rng.Uint32N(span)
	}
	return out
}

func sortedUnique(values []uint32) []uint32 {
	out := slices.Clone(values)
	slices.Sort(out)
	out = slices.Compact(out)
	if out == nil {
		return []uint32{}
	}
	return out
}

func setOf(values []uint32) map[uint32]struct{} {
	set := make(map[uint32]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
}

func intersectRef(a, b []uint32) []uint32 {
	set := setOf(b)
	out := []uint32{}
	for _, v := range a {
		if _, ok := set[v]; ok {
			out = append(out, v)
		}
	}
	return out
}

func unionRef(a, b []uint32) []uint32 {
	return sortedUnique(append(slices.Clone(a), b...))
}

func differenceRef(a, b []uint32) []uint32 {
	set := setOf(b)
	out := []uint32{}
	for _, v := range a {
		if _, ok := set[v]; !ok {
			out = append(out, v)
		}
	}
	return out
}

func xorRef(a, b []uint32) []uint32 {
	out := append(differenceRef(a, b), differenceRef(b, a)...)
	slices.Sort(out)
	return out
}
