package bench

import (
	"encoding/binary"
	"sort"
	"strconv"
	"testing"

	roaring "goalgo/roaring_bitmaps"
)

// caseName labels a benchmark case as key=value pairs, which is the shape
// benchstat -col and benchdraw need to slice the results by dimension.
func caseName(runopt, data string) string {
	return "runopt=" + runopt + "/data=" + data
}

// pair holds two bitmaps of the same variant and dataset, ready for a binary
// operation. Building them is not part of what the set-operation benchmarks
// measure, so it happens once per case.
type pair struct {
	a Bitmap
	b Bitmap
}

func buildPair(v variant, left, right dataset) pair {
	return pair{a: v.new(left.values...), b: v.new(right.values...)}
}

// eachCase runs fn once per variant and dataset.
func eachCase(b *testing.B, fn func(b *testing.B, v variant, left, right dataset)) {
	left := datasets(0)
	right := datasets(1 << 12)

	for _, v := range variants() {
		for i, ds := range left {
			b.Run(caseName(v.name, ds.name), func(b *testing.B) {
				fn(b, v, ds, right[i])
			})
		}
	}
}

// BenchmarkBuild measures the constructor; for the on variant it includes the
// RunOptimize pass.
func BenchmarkBuild(b *testing.B) {
	eachCase(b, func(b *testing.B, v variant, left, _ dataset) {
		b.ReportAllocs()
		for b.Loop() {
			_ = v.new(left.values...)
		}
	})
}

// BenchmarkMemory does no timing work worth reading; it reports the heap bytes
// the bitmap occupies per value, which is what run containers save. The set
// operations report their own allocation per call through -benchmem.
//
// zipf is measured at twice the usual size: its crowded head is where the
// container kinds mix, and doubling the draw makes that head dense enough to
// show how the footprint moves as chunks cross from arrays into bitmaps.
func BenchmarkMemory(b *testing.B) {
	sets := datasets(0)
	for _, ds := range datasetsOfSize(2*datasetSize(), 0) {
		if ds.name == "zipf" {
			sets[2] = ds
		}
	}

	for _, v := range variants() {
		for _, ds := range sets {
			b.Run(caseName(v.name, ds.name), func(b *testing.B) {
				bm := v.new(ds.values...)
				for b.Loop() {
					_ = bm.SizeInBytes()
				}
				// Reported after the loop so the timing harness does not clear it.
				b.ReportMetric(float64(bm.SizeInBytes())/float64(bm.Cardinality()), "bytes/value")
			})
		}
	}
}

func BenchmarkContains(b *testing.B) {
	eachCase(b, func(b *testing.B, v variant, left, _ dataset) {
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
	eachCase(b, func(b *testing.B, v variant, left, right dataset) {
		p := buildPair(v, left, right)

		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			op(p)
		}
	})
}

// Serialisation has no "off" variant to compare with.
// Its point is made against the alternative one would ship without it: the set
// kept as a sorted []uint32 beside the bitmap, written out as is and read back
// into the same kind of slice. Like the Rank/Select baseline, that copy is built
// once in prepare, not on every call. Cases are tagged codec=raw and codec=roaring.

type codec struct {
	name string
	// prepare builds whatever the codec needs beside the bitmap. It returns the
	// encoder, a decoder that only reaches the state the codec's user works
	// with — a slice for raw, a bitmap for roaring — and a reader that lists the
	// values of the last decoded state, used by the agreement test alone.
	prepare func(bm *roaring.Bitmap) (encode func() []byte, decode func([]byte) error, decoded func() []uint32)
}

func codecs() []codec {
	return []codec{
		{
			name: "raw",
			prepare: func(bm *roaring.Bitmap) (func() []byte, func([]byte) error, func() []uint32) {
				values := bm.ToArray()
				var last []uint32
				encode := func() []byte {
					out := make([]byte, 0, 4*len(values))
					for _, v := range values {
						out = binary.LittleEndian.AppendUint32(out, v)
					}
					return out
				}
				decode := func(data []byte) error {
					last = make([]uint32, len(data)/4)
					for i := range last {
						last[i] = binary.LittleEndian.Uint32(data[4*i:])
					}
					return nil
				}
				return encode, decode, func() []uint32 { return last }
			},
		},
		{
			name: "roaring",
			prepare: func(bm *roaring.Bitmap) (func() []byte, func([]byte) error, func() []uint32) {
				var last *roaring.Bitmap
				decode := func(data []byte) error {
					got, err := roaring.FromBytes(data)
					last = got
					return err
				}
				return bm.ToBytes, decode, func() []uint32 { return last.ToArray() }
			},
		},
	}
}

type codecFns struct {
	encode  func() []byte
	decode  func([]byte) error
	decoded func() []uint32
}

func eachCodec(b *testing.B, fn func(b *testing.B, c codecFns, bm *roaring.Bitmap, data []byte)) {
	for _, ds := range datasets(0) {
		bm := roaring.New(ds.values...)
		bm.RunOptimize()
		for _, c := range codecs() {
			b.Run("codec="+c.name+"/data="+ds.name, func(b *testing.B) {
				encode, decode, decoded := c.prepare(bm)
				fn(b, codecFns{encode, decode, decoded}, bm, encode())
			})
		}
	}
}

// BenchmarkSerializedSize does no timing work worth reading; it reports the bytes
// each codec needs per value, which is the other half of the comparison.
func BenchmarkSerializedSize(b *testing.B) {
	eachCodec(b, func(b *testing.B, _ codecFns, bm *roaring.Bitmap, data []byte) {
		for b.Loop() {
			_ = len(data)
		}
		// Reported after the loop so the timing harness does not clear it.
		b.ReportMetric(float64(len(data))/float64(bm.Cardinality()), "bytes/value")
	})
}

func BenchmarkToBytes(b *testing.B) {
	eachCodec(b, func(b *testing.B, c codecFns, _ *roaring.Bitmap, _ []byte) {
		b.ReportAllocs()
		for b.Loop() {
			_ = c.encode()
		}
	})
}

// BenchmarkFromBytes measures decoding only: roaring lands in containers, raw lands
// in a slice, and neither converts its result into anything else.
func BenchmarkFromBytes(b *testing.B) {
	eachCodec(b, func(b *testing.B, c codecFns, _ *roaring.Bitmap, data []byte) {
		b.ReportAllocs()
		for b.Loop() {
			if err := c.decode(data); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Rank and Select have no "off" variant either, so like serialisation they are compared
// against what one would do without them: keep a sorted copy of the values next
// to the bitmap and answer from it — a binary search for Rank, an index for
// Select. That copy is fast but costs 4 bytes per value on top of the bitmap;
// roaring answers from the containers plus a prefix-sum cache of one int per
// container. Cases are tagged method=array and method=roaring, and each reports the
// bytes of auxiliary state per value so the trade-off is on the chart.

type positional struct {
	name string
	// prepare builds whatever the method needs beside the bitmap and returns
	// the two query functions plus the size of that state in bytes.
	prepare func(bm *roaring.Bitmap) (rank func(uint32) int, sel func(int) uint32, auxBytes int)
}

func positionals() []positional {
	return []positional{
		{
			name: "array",
			prepare: func(bm *roaring.Bitmap) (func(uint32) int, func(int) uint32, int) {
				values := bm.ToArray()
				rank := func(v uint32) int {
					return sort.Search(len(values), func(i int) bool { return values[i] > v })
				}
				sel := func(i int) uint32 { return values[i] }
				return rank, sel, 4 * len(values)
			},
		},
		{
			name: "roaring",
			prepare: func(bm *roaring.Bitmap) (func(uint32) int, func(int) uint32, int) {
				sel := func(i int) uint32 {
					v, _ := bm.Select(i)
					return v
				}
				// The prefix-sum cache: one int per container plus a sentinel.
				arrays, bitmaps, runs := bm.Stats()
				return bm.Rank, sel, 8 * (arrays + bitmaps + runs + 1)
			},
		},
	}
}

func eachPositional(b *testing.B, fn func(b *testing.B, rank func(uint32) int, sel func(int) uint32, bm *roaring.Bitmap, values []uint32)) {
	for _, ds := range datasets(0) {
		bm := roaring.New(ds.values...)
		bm.RunOptimize()
		for _, m := range positionals() {
			b.Run("method="+m.name+"/data="+ds.name, func(b *testing.B) {
				rank, sel, aux := m.prepare(bm)
				b.ReportAllocs()
				b.ResetTimer()
				fn(b, rank, sel, bm, ds.values)
				b.ReportMetric(float64(aux)/float64(bm.Cardinality()), "aux-bytes/value")
			})
		}
	}
}

// BenchmarkRank asks for the rank of values drawn from the set itself.
func BenchmarkRank(b *testing.B) {
	eachPositional(b, func(b *testing.B, rank func(uint32) int, _ func(int) uint32, _ *roaring.Bitmap, values []uint32) {
		i, sum := 0, 0
		for b.Loop() {
			sum += rank(values[i])
			if i++; i == len(values) {
				i = 0
			}
		}
		if sum == 0 {
			b.Fatal("ranks summed to zero: the benchmark is not doing the work it claims")
		}
	})
}

// BenchmarkSelect asks for positions spread across the whole set.
func BenchmarkSelect(b *testing.B) {
	eachPositional(b, func(b *testing.B, _ func(uint32) int, sel func(int) uint32, bm *roaring.Bitmap, _ []uint32) {
		n := bm.Cardinality()
		i, sum := 0, uint64(0)
		for b.Loop() {
			sum += uint64(sel(i))
			if i += 7919; i >= n {
				i -= n
			}
		}
		if sum == 0 {
			b.Fatal("selected values summed to zero: the benchmark is not doing the work it claims")
		}
	})
}

// TestPositionalsAgree checks that both methods answer the same.
func TestPositionalsAgree(t *testing.T) {
	t.Parallel()

	for _, ds := range datasets(0) {
		t.Run(ds.name, func(t *testing.T) {
			t.Parallel()

			bm := roaring.New(ds.values...)
			bm.RunOptimize()
			ms := positionals()
			rankA, selA, _ := ms[0].prepare(bm)
			rankB, selB, _ := ms[1].prepare(bm)
			n := bm.Cardinality()
			for i := 0; i < n; i += 997 {
				if a, b := selA(i), selB(i); a != b {
					t.Fatalf("Select(%d): array %d, roaring %d", i, a, b)
				}
			}
			for _, v := range ds.values[:200] {
				if a, b := rankA(v), rankB(v); a != b {
					t.Fatalf("Rank(%d): array %d, roaring %d", v, a, b)
				}
			}
		})
	}
}

// TestVariantsAgree guards the comparison itself: measuring RunOptimize is only
// meaningful while the bitmap answers the same with and without it.
func TestVariantsAgree(t *testing.T) {
	t.Parallel()

	left := datasets(0)
	right := datasets(1 << 12)

	for i, ds := range left {
		t.Run(ds.name, func(t *testing.T) {
			t.Parallel()

			var want []string
			for _, v := range variants() {
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
						t.Fatalf("runopt=%s disagrees on operation %d: %s vs %s", v.name, j, got[j], want[j])
					}
				}
			}
		})
	}
}

// TestCodecsAgree checks that both codecs round-trip to the same values.
func TestCodecsAgree(t *testing.T) {
	t.Parallel()

	for _, ds := range datasets(0) {
		t.Run(ds.name, func(t *testing.T) {
			t.Parallel()

			bm := roaring.New(ds.values...)
			bm.RunOptimize()
			want := bm.ToArray()

			for _, c := range codecs() {
				encode, decode, decoded := c.prepare(bm)
				if err := decode(encode()); err != nil {
					t.Fatalf("%s: %v", c.name, err)
				}
				got := decoded()
				if len(got) != len(want) {
					t.Fatalf("%s round-trip changed the cardinality: %d vs %d", c.name, len(got), len(want))
				}
				for k := range want {
					if got[k] != want[k] {
						t.Fatalf("%s round-trip changed value %d: %d vs %d", c.name, k, got[k], want[k])
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
	values := bm.(bitmap).ToArray()

	var sum uint64
	for _, v := range values {
		sum = sum*1099511628211 ^ uint64(v)
	}
	return strconv.Itoa(len(values)) + ":" + strconv.FormatUint(sum, 16)
}
