package bench

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	roaring "goalgo/roaring_bitmaps"
	"goalgo/roaring_bitmaps/internal/simdops"
)

// The SIMD benchmarks measure what the vector build of simdops buys: the same
// operations over the same datasets, run once on each build. A single binary
// carries one simdops implementation, so the comparison is two passes — the
// vector build and the scalar one — whose cases are tagged simd=on and simd=off
// and concatenated into one file (make roaring_bench_simd).
//
// RunOptimize is deliberately not applied. It turns the runs dataset into run
// containers, which stay in interval arithmetic and never touch a vector
// primitive; without it those chunks are bitmap containers, which is exactly
// where the fused AVX-512 loops run. sparse is all tiny array containers and
// zipf mixes the two, so the three datasets together show where vectorisation
// pays and where it cannot.

func simdTag() string {
	if simdops.Vectorized {
		return "on"
	}
	return "off"
}

func simdCaseName(data string) string {
	return "simd=" + simdTag() + "/data=" + data
}

// eachSIMDCase runs fn once per dataset on the plain bitmap of this build.
func eachSIMDCase(b *testing.B, fn func(b *testing.B, left, right dataset)) {
	left := datasets(0)
	right := datasets(1 << 12)

	for i, ds := range left {
		b.Run(simdCaseName(ds.name), func(b *testing.B) {
			fn(b, ds, right[i])
		})
	}
}

func BenchmarkSIMDContains(b *testing.B) {
	eachSIMDCase(b, func(b *testing.B, left, _ dataset) {
		bm := roaring.New(left.values...)
		values := left.values

		b.ResetTimer()
		i, hits := 0, 0
		for b.Loop() {
			if bm.Contains(values[i]) {
				hits++
			}
			if i++; i == len(values) {
				i = 0
			}
		}
		if hits == 0 {
			b.Fatal("no hits: the benchmark is not doing the work it claims")
		}
	})
}

func BenchmarkSIMDAnd(b *testing.B) {
	benchmarkSIMDBinary(b, func(a, c *roaring.Bitmap) { _ = a.And(c) })
}

func BenchmarkSIMDOr(b *testing.B) {
	benchmarkSIMDBinary(b, func(a, c *roaring.Bitmap) { _ = a.Or(c) })
}

func BenchmarkSIMDAndNot(b *testing.B) {
	benchmarkSIMDBinary(b, func(a, c *roaring.Bitmap) { _ = a.AndNot(c) })
}

func BenchmarkSIMDXor(b *testing.B) {
	benchmarkSIMDBinary(b, func(a, c *roaring.Bitmap) { _ = a.Xor(c) })
}

func BenchmarkSIMDAndCardinality(b *testing.B) {
	benchmarkSIMDBinary(b, func(a, c *roaring.Bitmap) { _ = a.AndCardinality(c) })
}

func benchmarkSIMDBinary(b *testing.B, op func(a, c *roaring.Bitmap)) {
	eachSIMDCase(b, func(b *testing.B, left, right dataset) {
		a := roaring.New(left.values...)
		c := roaring.New(right.values...)

		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			op(a, c)
		}
	})
}

// BenchmarkSIMDToArray materialises the whole set, which for bitmap containers
// is the mask-to-positions primitive.
func BenchmarkSIMDToArray(b *testing.B) {
	eachSIMDCase(b, func(b *testing.B, left, _ dataset) {
		bm := roaring.New(left.values...)

		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			_ = bm.ToArray()
		}
	})
}

// TestSIMDOperationsConsistent checks, on whichever simdops build is compiled
// in, that the measured operations agree with each other over the benchmark
// datasets: the set identities that tie And, Or, AndNot, Xor and
// AndCardinality together, and that ToArray lists exactly the set. Run under
// both builds it is the evidence that the two passes time the same answers.
func TestSIMDOperationsConsistent(t *testing.T) {
	t.Parallel()

	left := datasets(0)
	right := datasets(1 << 12)

	tests := []struct {
		name  string
		left  dataset
		right dataset
	}{
		{name: "sparse", left: left[0], right: right[0]},
		{name: "runs", left: left[1], right: right[1]},
		{name: "zipf", left: left[2], right: right[2]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := roaring.New(tt.left.values...)
			c := roaring.New(tt.right.values...)

			and := a.And(c)
			or := a.Or(c)
			andNot := a.AndNot(c)
			xor := a.Xor(c)

			require.Equal(t, and.Cardinality(), a.AndCardinality(c))
			require.Equal(t, a.Cardinality()+c.Cardinality(), or.Cardinality()+and.Cardinality())
			require.Equal(t, a.Cardinality()-and.Cardinality(), andNot.Cardinality())
			require.Equal(t, or.Cardinality()-and.Cardinality(), xor.Cardinality())

			values := a.ToArray()
			require.Len(t, values, a.Cardinality())
			require.True(t, slices.IsSorted(values))
			for _, v := range values[:min(len(values), 1000)] {
				require.True(t, a.Contains(v))
			}
			for _, v := range and.ToArray()[:min(and.Cardinality(), 1000)] {
				require.True(t, a.Contains(v) && c.Contains(v))
			}
		})
	}
}
