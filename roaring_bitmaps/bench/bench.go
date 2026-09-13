// Package bench drives v1, v2 and v3 through the same workloads so the three
// container models can be compared directly.
//
// The versions have identical public APIs but distinct types, so each is wrapped
// in a thin adapter behind one interface. The wrappers only forward calls; any
// difference the benchmarks show comes from the container model, not from here.
package bench

import (
	"math/rand/v2"
	"os"
	"strconv"

	roaringv1 "goalgo/roaring_bitmaps/v1"
	roaringv2 "goalgo/roaring_bitmaps/v2"
	roaringv3 "goalgo/roaring_bitmaps/v3"
)

// Bitmap is the slice of the API every version shares.
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

type v1Bitmap struct{ *roaringv1.Bitmap }

func (b v1Bitmap) And(other Bitmap) Bitmap {
	return v1Bitmap{b.Bitmap.And(other.(v1Bitmap).Bitmap)}
}
func (b v1Bitmap) Or(other Bitmap) Bitmap { return v1Bitmap{b.Bitmap.Or(other.(v1Bitmap).Bitmap)} }
func (b v1Bitmap) AndNot(other Bitmap) Bitmap {
	return v1Bitmap{b.Bitmap.AndNot(other.(v1Bitmap).Bitmap)}
}
func (b v1Bitmap) Xor(other Bitmap) Bitmap { return v1Bitmap{b.Bitmap.Xor(other.(v1Bitmap).Bitmap)} }
func (b v1Bitmap) AndCardinality(other Bitmap) int {
	return b.Bitmap.AndCardinality(other.(v1Bitmap).Bitmap)
}

type v2Bitmap struct{ *roaringv2.Bitmap }

func (b v2Bitmap) And(other Bitmap) Bitmap {
	return v2Bitmap{b.Bitmap.And(other.(v2Bitmap).Bitmap)}
}
func (b v2Bitmap) Or(other Bitmap) Bitmap { return v2Bitmap{b.Bitmap.Or(other.(v2Bitmap).Bitmap)} }
func (b v2Bitmap) AndNot(other Bitmap) Bitmap {
	return v2Bitmap{b.Bitmap.AndNot(other.(v2Bitmap).Bitmap)}
}
func (b v2Bitmap) Xor(other Bitmap) Bitmap { return v2Bitmap{b.Bitmap.Xor(other.(v2Bitmap).Bitmap)} }
func (b v2Bitmap) AndCardinality(other Bitmap) int {
	return b.Bitmap.AndCardinality(other.(v2Bitmap).Bitmap)
}

type v3Bitmap struct{ *roaringv3.Bitmap }

func (b v3Bitmap) And(other Bitmap) Bitmap {
	return v3Bitmap{b.Bitmap.And(other.(v3Bitmap).Bitmap)}
}
func (b v3Bitmap) Or(other Bitmap) Bitmap { return v3Bitmap{b.Bitmap.Or(other.(v3Bitmap).Bitmap)} }
func (b v3Bitmap) AndNot(other Bitmap) Bitmap {
	return v3Bitmap{b.Bitmap.AndNot(other.(v3Bitmap).Bitmap)}
}
func (b v3Bitmap) Xor(other Bitmap) Bitmap { return v3Bitmap{b.Bitmap.Xor(other.(v3Bitmap).Bitmap)} }
func (b v3Bitmap) AndCardinality(other Bitmap) int {
	return b.Bitmap.AndCardinality(other.(v3Bitmap).Bitmap)
}

type version struct {
	name string
	new  func(values ...uint32) Bitmap
}

// versions lists the implementations under comparison, each built the way it is
// meant to be used. v2 and v3 exist for their run containers, so their
// constructors include RunOptimize: the difference between v1 and v2 in any chart
// is exactly what that encoding buys, or costs.
func versions() []version {
	return []version{
		{name: "v1", new: func(values ...uint32) Bitmap {
			return v1Bitmap{roaringv1.New(values...)}
		}},
		{name: "v2", new: func(values ...uint32) Bitmap {
			b := roaringv2.New(values...)
			b.RunOptimize()
			return v2Bitmap{b}
		}},
		{name: "v3", new: func(values ...uint32) Bitmap {
			b := roaringv3.New(values...)
			b.RunOptimize()
			return v3Bitmap{b}
		}},
	}
}

type dataset struct {
	name   string
	values []uint32
}

// datasetSize is the number of values per dataset. Override it with ROARING_N to
// see how the versions separate as the data grows.
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
//   - runs: stretches of 512 consecutive values separated by random gaps. In v1
//     these fill bitmap containers; in v2 and v3 they collapse into a few
//     intervals each. This is the shape run containers were built for.
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
