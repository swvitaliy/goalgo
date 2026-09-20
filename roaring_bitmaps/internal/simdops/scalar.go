//go:build !(goexperiment.simd && amd64)

package simdops

import (
	"math/bits"
	"slices"
)

// Vectorized reports that this build uses the plain-Go primitives.
const Vectorized = false

// AndTo stores a&b into dst and returns the number of set bits written.
// dst may alias a or b. All three slices must be the same length.
func AndTo(dst, a, b []uint64) int {
	n := checkTriple(dst, a, b)
	total := 0
	for i := range n {
		dst[i] = a[i] & b[i]
		total += bits.OnesCount64(dst[i])
	}
	return total
}

// OrTo stores a|b into dst and returns the number of set bits written.
func OrTo(dst, a, b []uint64) int {
	n := checkTriple(dst, a, b)
	total := 0
	for i := range n {
		dst[i] = a[i] | b[i]
		total += bits.OnesCount64(dst[i])
	}
	return total
}

// AndNotTo stores a&^b into dst and returns the number of set bits written.
func AndNotTo(dst, a, b []uint64) int {
	n := checkTriple(dst, a, b)
	total := 0
	for i := 0; i < n; i++ {
		dst[i] = a[i] &^ b[i]
		total += bits.OnesCount64(dst[i])
	}
	return total
}

// XorTo stores a^b into dst and returns the number of set bits written.
func XorTo(dst, a, b []uint64) int {
	n := checkTriple(dst, a, b)
	total := 0
	for i := range n {
		dst[i] = a[i] ^ b[i]
		total += bits.OnesCount64(dst[i])
	}
	return total
}

// AndCardinality returns the number of set bits in a&b without materialising the
// result, for callers that only need the count.
func AndCardinality(a, b []uint64) int {
	total := 0
	for i := range min(len(a), len(b)) {
		total += bits.OnesCount64(a[i] & b[i])
	}
	return total
}

// OrCardinality returns the number of set bits in a|b without materialising it.
func OrCardinality(a, b []uint64) int {
	total := 0
	for i := range min(len(a), len(b)) {
		total += bits.OnesCount64(a[i] | b[i])
	}
	return total
}

// AndNotCardinality returns the number of set bits in a&^b without materialising it.
func AndNotCardinality(a, b []uint64) int {
	total := 0
	for i := range min(len(a), len(b)) {
		total += bits.OnesCount64(a[i] &^ b[i])
	}
	return total
}

// XorCardinality returns the number of set bits in a^b without materialising it.
func XorCardinality(a, b []uint64) int {
	total := 0
	for i := range min(len(a), len(b)) {
		total += bits.OnesCount64(a[i] ^ b[i])
	}
	return total
}

// Popcount returns the number of set bits in a.
func Popcount(a []uint64) int {
	total := 0
	for _, w := range a {
		total += bits.OnesCount64(w)
	}
	return total
}

// Intersects reports whether a&b has any set bit, stopping at the first hit.
func Intersects(a, b []uint64) bool {
	for i := range min(len(a), len(b)) {
		if a[i]&b[i] != 0 {
			return true
		}
	}
	return false
}

// BitmapToArray writes the positions of the set bits of bm, in ascending order,
// into dst and returns how many were written. dst must hold Popcount(bm) values.
func BitmapToArray(dst []uint16, bm []uint64) int {
	n := 0
	for wi, w := range bm {
		base := wi * 64
		for w != 0 {
			dst[n] = uint16(base + bits.TrailingZeros64(w))
			n++
			w &= w - 1
		}
	}
	return n
}

// IntersectArrays writes the intersection of the sorted, duplicate-free arrays a
// and b into dst and returns how many values were written. dst must have room for
// min(len(a), len(b)) values and must not alias a or b.
func IntersectArrays(dst, a, b []uint16) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	small, large := a, b
	if len(small) > len(large) {
		small, large = large, small
	}
	if len(large)/len(small) >= gallopRatio {
		return intersectGallop(dst, small, large)
	}
	return intersectLinear(dst, a, b)
}

// intersectGallop looks each value of the small array up in the large one with a
// binary search that only ever moves forward, since both arrays are sorted.
func intersectGallop(dst, small, large []uint16) int {
	n, lo := 0, 0
	for _, x := range small {
		i, found := slices.BinarySearch(large[lo:], x)
		lo += i
		if lo >= len(large) {
			break
		}
		if found {
			dst[n] = x
			n++
		}
	}
	return n
}

func fillWords(dst []uint64, v uint64) {
	for i := range dst {
		dst[i] = v
	}
}

func flipWords(dst []uint64) {
	for i := range dst {
		dst[i] = ^dst[i]
	}
}

func anyNonZero(words []uint64) bool {
	for _, w := range words {
		if w != 0 {
			return true
		}
	}
	return false
}
