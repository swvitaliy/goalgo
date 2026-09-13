//go:build goexperiment.simd && amd64

package bench

import (
	"encoding/binary"
	"strconv"
	"testing"

	roaringv3 "goalgo/roaring_bitmaps/v3"
)

// caseName labels a benchmark case as key=value pairs, which is the shape
// benchstat -col and benchdraw need to slice the results by dimension.
func caseName(version, data string) string {
	return "version=" + version + "/data=" + data
}

// pair holds two bitmaps of the same version and dataset, ready for a binary
// operation. Building them is not part of what the set-operation benchmarks
// measure, so it happens once per case.
type pair struct {
	a Bitmap
	b Bitmap
}

func buildPair(v version, left, right dataset) pair {
	return pair{a: v.new(left.values...), b: v.new(right.values...)}
}

// eachCase runs fn once per version and dataset.
func eachCase(b *testing.B, fn func(b *testing.B, v version, left, right dataset)) {
	left := datasets(0)
	right := datasets(1 << 12)

	for _, v := range versions() {
		for i, ds := range left {
			b.Run(caseName(v.name, ds.name), func(b *testing.B) {
				fn(b, v, ds, right[i])
			})
		}
	}
}

// BenchmarkBuild measures the constructor as each version is meant to be used,
// so for v2 and v3 it includes the RunOptimize pass.
func BenchmarkBuild(b *testing.B) {
	eachCase(b, func(b *testing.B, v version, left, _ dataset) {
		b.ReportAllocs()
		for b.Loop() {
			_ = v.new(left.values...)
		}
	})
}

func BenchmarkContains(b *testing.B) {
	eachCase(b, func(b *testing.B, v version, left, _ dataset) {
		bm := v.new(left.values...)
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
	eachCase(b, func(b *testing.B, v version, left, right dataset) {
		p := buildPair(v, left, right)

		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			op(p)
		}
	})
}

// Serialisation exists only in v3, so there is no version to compare it with.
// Its point is made against the alternative one would ship without it: a dump of
// the sorted uint32 values, rebuilt into a bitmap on the other side. Cases are
// tagged codec=raw and codec=v3.

type codec struct {
	name   string
	encode func(bm *roaringv3.Bitmap) []byte
	decode func(data []byte) (*roaringv3.Bitmap, error)
}

func codecs() []codec {
	return []codec{
		{
			name: "raw",
			encode: func(bm *roaringv3.Bitmap) []byte {
				values := bm.ToArray()
				out := make([]byte, 0, 4*len(values))
				for _, v := range values {
					out = binary.LittleEndian.AppendUint32(out, v)
				}
				return out
			},
			decode: func(data []byte) (*roaringv3.Bitmap, error) {
				values := make([]uint32, len(data)/4)
				for i := range values {
					values[i] = binary.LittleEndian.Uint32(data[4*i:])
				}
				// Reach the same in-memory state FromBytes produces from a
				// run-optimised stream, so the two decoders end at the same place.
				bm := roaringv3.New(values...)
				bm.RunOptimize()
				return bm, nil
			},
		},
		{
			name:   "v3",
			encode: (*roaringv3.Bitmap).ToBytes,
			decode: roaringv3.FromBytes,
		},
	}
}

func eachCodec(b *testing.B, fn func(b *testing.B, c codec, bm *roaringv3.Bitmap, data []byte)) {
	for _, ds := range datasets(0) {
		bm := roaringv3.New(ds.values...)
		bm.RunOptimize()
		for _, c := range codecs() {
			b.Run("codec="+c.name+"/data="+ds.name, func(b *testing.B) {
				fn(b, c, bm, c.encode(bm))
			})
		}
	}
}

// BenchmarkSerializedSize does no timing work worth reading; it reports the bytes
// each codec needs per value, which is the other half of the comparison.
func BenchmarkSerializedSize(b *testing.B) {
	eachCodec(b, func(b *testing.B, _ codec, bm *roaringv3.Bitmap, data []byte) {
		for b.Loop() {
			_ = len(data)
		}
		// Reported after the loop so the timing harness does not clear it.
		b.ReportMetric(float64(len(data))/float64(bm.Cardinality()), "bytes/value")
	})
}

func BenchmarkToBytes(b *testing.B) {
	eachCodec(b, func(b *testing.B, c codec, bm *roaringv3.Bitmap, _ []byte) {
		b.ReportAllocs()
		for b.Loop() {
			_ = c.encode(bm)
		}
	})
}

func BenchmarkFromBytes(b *testing.B) {
	eachCodec(b, func(b *testing.B, c codec, _ *roaringv3.Bitmap, data []byte) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := c.decode(data); err != nil {
				b.Fatal(err)
			}
		}
	})
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
				p := buildPair(v, ds, right[i])
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

// TestCodecsAgree checks that both codecs round-trip to the same set.
func TestCodecsAgree(t *testing.T) {
	t.Parallel()

	for _, ds := range datasets(0) {
		t.Run(ds.name, func(t *testing.T) {
			t.Parallel()

			bm := roaringv3.New(ds.values...)
			bm.RunOptimize()
			want := summarize(v3Bitmap{bm})

			for _, c := range codecs() {
				got, err := c.decode(c.encode(bm))
				if err != nil {
					t.Fatalf("%s: %v", c.name, err)
				}
				if s := summarize(v3Bitmap{got}); s != want {
					t.Fatalf("%s round-trip changed the set: %s vs %s", c.name, s, want)
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
