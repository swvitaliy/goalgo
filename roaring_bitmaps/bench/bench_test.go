//go:build goexperiment.simd && amd64

package bench

import (
	"strconv"
	"testing"

	roaringv3 "goalgo/roaring_bitmaps/v3"
)

// pair holds two bitmaps of the same version and dataset, ready for a binary
// operation. Building them is not part of what the set-operation benchmarks
// measure, so it happens once per case.
type pair struct {
	a Bitmap
	b Bitmap
}

func buildPair(v version, left, right dataset, optimize bool) pair {
	p := pair{a: v.new(left.values...), b: v.new(right.values...)}
	if optimize {
		if o, ok := p.a.(runOptimizer); ok {
			o.RunOptimize()
		}
		if o, ok := p.b.(runOptimizer); ok {
			o.RunOptimize()
		}
	}
	return p
}

// eachCase runs fn for every version, dataset and run-encoding setting. v1 has no
// run containers, so it is only visited once per dataset.
func eachCase(b *testing.B, fn func(b *testing.B, v version, left, right dataset, optimize bool)) {
	left := datasets(0)
	right := datasets(1 << 12)

	for _, v := range versions() {
		for i, ds := range left {
			for _, optimize := range []bool{false, true} {
				name := v.name + "/" + ds.name
				if optimize {
					if v.name == "v1" {
						continue
					}
					name += "/runopt"
				}
				b.Run(name, func(b *testing.B) {
					fn(b, v, ds, right[i], optimize)
				})
			}
		}
	}
}

func BenchmarkBuild(b *testing.B) {
	eachCase(b, func(b *testing.B, v version, left, _ dataset, optimize bool) {
		b.ReportAllocs()
		for b.Loop() {
			bm := v.new(left.values...)
			if optimize {
				bm.(runOptimizer).RunOptimize()
			}
		}
	})
}

func BenchmarkContains(b *testing.B) {
	eachCase(b, func(b *testing.B, v version, left, _ dataset, optimize bool) {
		p := buildPair(v, left, left, optimize)
		values := left.values

		b.ResetTimer()
		i, hits := 0, 0
		for b.Loop() {
			if p.a.Contains(values[i]) {
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

func BenchmarkAnd(b *testing.B) {
	benchmarkBinary(b, func(p pair) { _ = p.a.And(p.b) })
}

func BenchmarkOr(b *testing.B) {
	benchmarkBinary(b, func(p pair) { _ = p.a.Or(p.b) })
}

func BenchmarkAndNot(b *testing.B) {
	benchmarkBinary(b, func(p pair) { _ = p.a.AndNot(p.b) })
}

func BenchmarkXor(b *testing.B) {
	benchmarkBinary(b, func(p pair) { _ = p.a.Xor(p.b) })
}

func BenchmarkAndCardinality(b *testing.B) {
	benchmarkBinary(b, func(p pair) { _ = p.a.AndCardinality(p.b) })
}

func benchmarkBinary(b *testing.B, op func(pair)) {
	eachCase(b, func(b *testing.B, v version, left, right dataset, optimize bool) {
		p := buildPair(v, left, right, optimize)

		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			op(p)
		}
	})
}

// BenchmarkSerializedSize does no timing work worth reading; it reports the bytes
// each encoding needs per value, which is the other half of the v1/v2/v3
// comparison.
func BenchmarkSerializedSize(b *testing.B) {
	for _, ds := range datasets(0) {
		for _, optimize := range []bool{false, true} {
			name := "v3/" + ds.name
			if optimize {
				name += "/runopt"
			}
			b.Run(name, func(b *testing.B) {
				bm := roaringv3.New(ds.values...)
				if optimize {
					bm.RunOptimize()
				}
				for b.Loop() {
					_ = bm.SerializedSize()
				}
				// Reported after the loop so the timing harness does not clear it.
				b.ReportMetric(float64(bm.SerializedSize())/float64(bm.Cardinality()), "bytes/value")
			})
		}
	}
}

func BenchmarkToBytes(b *testing.B) {
	benchmarkSerialization(b, func(b *testing.B, bm *roaringv3.Bitmap, _ []byte) {
		for b.Loop() {
			_ = bm.ToBytes()
		}
	})
}

func BenchmarkFromBytes(b *testing.B) {
	benchmarkSerialization(b, func(b *testing.B, _ *roaringv3.Bitmap, data []byte) {
		for b.Loop() {
			if _, err := roaringv3.FromBytes(data); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkSerialization(b *testing.B, fn func(b *testing.B, bm *roaringv3.Bitmap, data []byte)) {
	for _, ds := range datasets(0) {
		for _, optimize := range []bool{false, true} {
			name := "v3/" + ds.name
			if optimize {
				name += "/runopt"
			}
			b.Run(name, func(b *testing.B) {
				bm := roaringv3.New(ds.values...)
				if optimize {
					bm.RunOptimize()
				}
				data := bm.ToBytes()

				b.SetBytes(int64(len(data)))
				b.ReportAllocs()
				b.ResetTimer()
				fn(b, bm, data)
			})
		}
	}
}

// TestVersionsAgree guards the comparison itself: a benchmark of three
// implementations is only meaningful while all three produce the same sets.
func TestVersionsAgree(t *testing.T) {
	t.Parallel()

	left := datasets(0)
	right := datasets(1 << 12)

	for i, ds := range left {
		t.Run(ds.name, func(t *testing.T) {
			t.Parallel()

			var want []string
			for _, v := range versions() {
				p := buildPair(v, ds, right[i], v.name != "v1")
				got := []string{
					summarize(p.a.And(p.b)),
					summarize(p.a.Or(p.b)),
					summarize(p.a.AndNot(p.b)),
					summarize(p.a.Xor(p.b)),
				}
				if want == nil {
					want = got
					continue
				}
				for j := range want {
					if got[j] != want[j] {
						t.Fatalf("%s disagrees with v1 on operation %d: %s vs %s", v.name, j, got[j], want[j])
					}
				}
			}
		})
	}
}

// summarize reduces a bitmap to a short fingerprint: cardinality plus a checksum
// over its values, which is enough to catch a disagreement without comparing
// millions of elements.
func summarize(bm Bitmap) string {
	var values []uint32
	switch v := bm.(type) {
	case v1Bitmap:
		values = v.ToArray()
	case v2Bitmap:
		values = v.ToArray()
	case v3Bitmap:
		values = v.ToArray()
	}

	var sum uint64
	for _, v := range values {
		sum = sum*1099511628211 ^ uint64(v)
	}
	return strconv.Itoa(len(values)) + ":" + strconv.FormatUint(sum, 16)
}
