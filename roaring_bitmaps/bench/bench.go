// Package bench drives the bitmap through the same workloads with and without
// RunOptimize, so the value of run containers can be read off directly.
//
// The bitmap is wrapped in a thin adapter behind one interface so the binary
// operations can be written once; the wrapper only forwards calls.
package bench

import (
	"math/rand/v2"
	"os"
	"strconv"

	roaring "goalgo/roaring_bitmaps"
)

// Bitmap is the slice of the API the benchmarks use.
type Bitmap interface {
	Add(v uint32) bool
	Contains(v uint32) bool
	Cardinality() int
	And(other Bitmap) Bitmap
	Or(other Bitmap) Bitmap
	AndNot(other Bitmap) Bitmap
	Xor(other Bitmap) Bitmap
	AndCardinality(other Bitmap) int
}

type bitmap struct{ *roaring.Bitmap }

func (b bitmap) And(other Bitmap) Bitmap {
	return bitmap{b.Bitmap.And(other.(bitmap).Bitmap)}
}
func (b bitmap) Or(other Bitmap) Bitmap { return bitmap{b.Bitmap.Or(other.(bitmap).Bitmap)} }
func (b bitmap) AndNot(other Bitmap) Bitmap {
	return bitmap{b.Bitmap.AndNot(other.(bitmap).Bitmap)}
}
func (b bitmap) Xor(other Bitmap) Bitmap { return bitmap{b.Bitmap.Xor(other.(bitmap).Bitmap)} }
func (b bitmap) AndCardinality(other Bitmap) int {
	return b.Bitmap.AndCardinality(other.(bitmap).Bitmap)
}

// variant is one way of using the bitmap. There is one implementation; what
// the benchmarks compare is whether RunOptimize was applied, which is the only
// thing that ever produces run containers. The difference between the two
// variants in any chart is exactly what that encoding buys, or costs.
type variant struct {
	name string
	new  func(values ...uint32) Bitmap
}

func variants() []variant {
	return []variant{
		{name: "off", new: func(values ...uint32) Bitmap {
			return bitmap{roaring.New(values...)}
		}},
		{name: "on", new: func(values ...uint32) Bitmap {
			b := roaring.New(values...)
			b.RunOptimize()
			return bitmap{b}
		}},
	}
}

type dataset struct {
	name   string
	values []uint32
}

// datasetSize is the number of values per dataset. Override it with ROARING_N to
// see how the variants separate as the data grows.
func datasetSize() int {
	if s := os.Getenv("ROARING_N"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			return n
		}
	}
	return 200_000
}

// datasets returns the three shapes that tell the container models apart.
//
//   - sparse: random values across 2^31. Every chunk holds a handful of values, so
//     everything is a tiny array container and per-container overhead dominates.
//   - runs: stretches of 512 consecutive values separated by random gaps.
//     Without RunOptimize these fill bitmap containers; with it they collapse
//     into a few intervals each. This is the shape run containers were built for.
//   - zipf: a skewed distribution with a crowded low end and a sparse tail, so one
//     bitmap mixes bitmap, array and run containers the way real data does.
func datasets(offset uint32) []dataset {
	n := datasetSize()
	rng := rand.New(rand.NewPCG(uint64(offset)+1, 0x5eed))

	sparse := make([]uint32, n)
	for i := range sparse {
		sparse[i] = offset + rng.Uint32N(1<<31)
	}

	runs := make([]uint32, 0, n)
	for v := offset; len(runs) < n; v += uint32(rng.UintN(4000) + 1) {
		for j := 0; j < 512 && len(runs) < n; j++ {
			runs = append(runs, v+uint32(j))
		}
	}

	// A skewed draw: most values land in a small prefix of the key space, a few
	// stray far out.
	zipf := make([]uint32, n)
	for i := range zipf {
		zipf[i] = offset + uint32(float64(1<<28)*rng.Float64()*rng.Float64()*rng.Float64())
	}

	return []dataset{
		{name: "sparse", values: sparse},
		{name: "runs", values: runs},
		{name: "zipf", values: zipf},
	}
}
