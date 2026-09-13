package cms

import (
	"math"

	"github.com/cespare/xxhash/v2"
)

type CMS struct {
	data   [][]uint64
	width  uint64
	hashes []func([]byte) uint64
}

func MakeHash(seed uint64) func(data []byte) uint64 {
	return func(data []byte) uint64 {
		digest := xxhash.NewWithSeed(seed)
		_, _ = digest.Write(data)
		return digest.Sum64()
	}
}

func NewCMS(width uint64, hashes []func([]byte) uint64) *CMS {
	data := make([][]uint64, len(hashes))
	for i := range data {
		data[i] = make([]uint64, width)
	}
	return &CMS{
		data:   data,
		width:  width,
		hashes: hashes,
	}
}

func (cms *CMS) Add(data []byte) {
	for i, h := range cms.hashes {
		index := h(data) % cms.width
		cms.data[i][index]++
	}
}

func (cms *CMS) CountMin(data []byte) uint64 {
	minValue := ^uint64(0)
	for i, h := range cms.hashes {
		index := h(data) % cms.width
		if cms.data[i][index] < minValue {
			minValue = cms.data[i][index]
		}
	}
	return minValue
}

func (cms *CMS) Precision() float64 {
	depth := len(cms.hashes)
	if depth == 0 {
		return 0.0
	}
	return 1 - math.Pow(math.E, 1/(float64(cms.width)/float64(depth)))
}

func Deviation(cms *CMS, data []byte, actualCount uint64) float64 {
	estimatedCount := cms.CountMin(data)
	if actualCount == 0 {
		return 0.0
	}
	return float64(absDiff(estimatedCount, actualCount)) / float64(actualCount)
}

func MaxDeviation(cms *CMS, data [][]byte, actualCounts []uint64) float64 {
	maxDev := 0.0
	for i, d := range data {
		dev := Deviation(cms, d, actualCounts[i%len(actualCounts)])
		if dev > maxDev {
			maxDev = dev
		}
	}
	return maxDev
}

func absDiff(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}
