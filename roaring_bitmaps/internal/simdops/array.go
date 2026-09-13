//go:build goexperiment.simd && amd64

package simdops

import (
	"math/bits"

	"simd/archsimd"
)

// gallopRatio is the size ratio at which intersecting two sorted arrays switches
// from a block-wise merge to a vector-probed galloping search. Below it the two
// arrays are close enough in size that a merge touches each element about once;
// above it, searching costs log2 of the large array per element of the small one,
// which wins by a wide margin.
const gallopRatio = 16

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
	return intersectMerge(dst, a, b)
}

// intersectMerge walks both arrays in 32-value blocks. For each pair of blocks it
// tests all 32x32 value combinations with 32 broadcast comparisons, accumulating
// one bit per matching lane of a, then compresses those lanes straight into dst.
func intersectMerge(dst, a, b []uint16) int {
	var buf [words16]uint16
	n, i, j := 0, 0, 0

	for i+words16 <= len(a) && j+words16 <= len(b) {
		va := archsimd.LoadUint16x32(a[i:])
		matched := uint32(0)
		for k := 0; k < words16; k++ {
			matched |= va.Equal(archsimd.BroadcastUint16x32(b[j+k])).ToBits()
		}
		if matched != 0 {
			va.Compress(archsimd.Mask16x32FromBits(matched)).StoreArray(&buf)
			n += copy(dst[n:], buf[:bits.OnesCount32(matched)])
		}

		// Advance past whichever block ends first; on a tie both are exhausted.
		lastA, lastB := a[i+words16-1], b[j+words16-1]
		if lastA <= lastB {
			i += words16
		}
		if lastB <= lastA {
			j += words16
		}
	}

	return n + intersectScalar(dst[n:], a[i:], b[j:])
}

// intersectGallop probes the large array for each value of the small one: an
// exponential search narrows the range to at most one vector, which a single
// broadcast comparison then tests.
func intersectGallop(dst, small, large []uint16) int {
	n, lo := 0, 0
	for _, x := range small {
		if lo >= len(large) || large[len(large)-1] < x {
			break
		}

		// Exponential search for an upper bound, starting where the previous
		// value left off: small is sorted, so the window only moves forward.
		hi, step := lo, 1
		for hi < len(large) && large[hi] < x {
			lo = hi + 1
			hi += step
			step *= 2
		}
		// Make the range half-open and inclusive of the first index known to
		// hold a value >= x, so a hit at that index is not searched past.
		if hi = min(hi+1, len(large)); lo >= hi {
			continue
		}
		for hi-lo > words16 {
			mid := int(uint(lo+hi) >> 1)
			if large[mid] < x {
				lo = mid + 1
			} else {
				hi = mid + 1
			}
		}

		window, _ := archsimd.LoadUint16x32Part(large[lo:hi])
		// Lanes past the window hold padding, so mask them out before testing.
		valid := uint32(1)<<uint(hi-lo) - 1
		if hi-lo == words16 {
			valid = ^uint32(0)
		}
		if window.Equal(archsimd.BroadcastUint16x32(x)).ToBits()&valid != 0 {
			dst[n] = x
			n++
		}
	}
	return n
}

// intersectScalar handles the tails left over once either array has fewer than a
// full vector of values remaining.
func intersectScalar(dst, a, b []uint16) int {
	n, i, j := 0, 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default:
			dst[n] = a[i]
			n++
			i++
			j++
		}
	}
	return n
}

// UnionArrays, DifferenceArrays and XorArrays stay scalar on purpose: their cost
// is dominated by writing a result whose length is not known in advance, so the
// per-element branch a merge pays is not what limits them.

// UnionArrays writes the union of the sorted arrays a and b into dst and returns
// how many values were written. dst must have room for len(a)+len(b) values.
func UnionArrays(dst, a, b []uint16) int {
	n, i, j := 0, 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			dst[n] = a[i]
			i++
		case a[i] > b[j]:
			dst[n] = b[j]
			j++
		default:
			dst[n] = a[i]
			i++
			j++
		}
		n++
	}
	n += copy(dst[n:], a[i:])
	n += copy(dst[n:], b[j:])
	return n
}

// DifferenceArrays writes a\b into dst and returns how many values were written.
// dst must have room for len(a) values.
func DifferenceArrays(dst, a, b []uint16) int {
	n, i, j := 0, 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			dst[n] = a[i]
			n++
			i++
		case a[i] > b[j]:
			j++
		default:
			i++
			j++
		}
	}
	return n + copy(dst[n:], a[i:])
}

// XorArrays writes the symmetric difference of a and b into dst and returns how
// many values were written. dst must have room for len(a)+len(b) values.
func XorArrays(dst, a, b []uint16) int {
	n, i, j := 0, 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			dst[n] = a[i]
			n++
			i++
		case a[i] > b[j]:
			dst[n] = b[j]
			n++
			j++
		default:
			i++
			j++
		}
	}
	n += copy(dst[n:], a[i:])
	n += copy(dst[n:], b[j:])
	return n
}
