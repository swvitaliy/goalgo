// Package simdops holds the bit-level primitives shared by every Roaring bitmap
// version in this directory.
//
// Keeping them in one place is deliberate: v1, v2 and v3 differ only in their
// container model, so benchmarks comparing the three measure that model rather
// than accidental differences in the low-level code.
//
// There are two implementations of the hot primitives, selected at build time:
//
//   - vector_amd64.go, under GOEXPERIMENT=simd on amd64, uses simd/archsimd
//     (AVX-512 with VPOPCNTDQ and VBMI2) and panics at init if the CPU lacks
//     them;
//   - scalar.go, in every other build, is plain Go over uint64 words.
//
// Both export the same API and pass the same tests. [Vectorized] reports which
// one is compiled in, so a benchmark can label its results. The scalar build is
// the baseline that shows what the vectorised one buys.
package simdops

import (
	"fmt"
	"math/bits"
)

// BitmapWords is the number of 64-bit words in a bitmap container, covering the
// full 65536-value key space of a single container.
const BitmapWords = 1024

// gallopRatio is the size ratio at which intersecting two sorted arrays switches
// from a merge to a galloping search. Below it the two arrays are close enough
// in size that a merge touches each element about once; above it, searching
// costs log2 of the large array per element of the small one, which wins by a
// wide margin.
const gallopRatio = 16

func checkTriple(dst, a, b []uint64) int {
	if len(a) != len(b) || len(dst) != len(a) {
		panic(fmt.Sprintf("simdops: length mismatch dst=%d a=%d b=%d", len(dst), len(a), len(b)))
	}
	return len(dst)
}

// intersectLinear is the plain merge intersection of two sorted arrays. The
// vector build uses it for the tails a full block cannot cover; the scalar build
// uses it outright.
func intersectLinear(dst, a, b []uint16) int {
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

// The range helpers all take a half-open bit range [start, end). Whole words in
// the middle go through fillWords/flipWords/anyNonZero, which is the part each
// build implements its own way; the two partial words at the edges are masked.

// SetRange sets every bit in [start, end).
func SetRange(bm []uint64, start, end int) {
	if start >= end {
		return
	}
	sw, ew := start/64, (end-1)/64
	if sw == ew {
		bm[sw] |= wordMask(start-sw*64, end-sw*64)
		return
	}
	bm[sw] |= wordMask(start-sw*64, 64)
	fillWords(bm[sw+1:ew], ^uint64(0))
	bm[ew] |= wordMask(0, end-ew*64)
}

// ClearRange clears every bit in [start, end).
func ClearRange(bm []uint64, start, end int) {
	if start >= end {
		return
	}
	sw, ew := start/64, (end-1)/64
	if sw == ew {
		bm[sw] &^= wordMask(start-sw*64, end-sw*64)
		return
	}
	bm[sw] &^= wordMask(start-sw*64, 64)
	fillWords(bm[sw+1:ew], 0)
	bm[ew] &^= wordMask(0, end-ew*64)
}

// FlipRange inverts every bit in [start, end).
func FlipRange(bm []uint64, start, end int) {
	if start >= end {
		return
	}
	sw, ew := start/64, (end-1)/64
	if sw == ew {
		bm[sw] ^= wordMask(start-sw*64, end-sw*64)
		return
	}
	bm[sw] ^= wordMask(start-sw*64, 64)
	flipWords(bm[sw+1 : ew])
	bm[ew] ^= wordMask(0, end-ew*64)
}

// CopyRange copies the bits of src in [start, end) into dst, leaving the bits of
// dst outside the range untouched.
func CopyRange(dst, src []uint64, start, end int) {
	if start >= end {
		return
	}
	sw, ew := start/64, (end-1)/64
	if sw == ew {
		m := wordMask(start-sw*64, end-sw*64)
		dst[sw] = dst[sw]&^m | src[sw]&m
		return
	}
	m := wordMask(start-sw*64, 64)
	dst[sw] = dst[sw]&^m | src[sw]&m
	copy(dst[sw+1:ew], src[sw+1:ew])
	m = wordMask(0, end-ew*64)
	dst[ew] = dst[ew]&^m | src[ew]&m
}

// wordMask returns a mask with bits [lo, hi) of a single word set, where hi is in
// 1..64.
func wordMask(lo, hi int) uint64 {
	// The parentheses matter: <<, >> and & share a precedence level in Go and
	// would otherwise group left to right.
	return (^uint64(0) << uint(lo)) & (^uint64(0) >> uint(64-hi))
}

// HasBitInRange reports whether any bit of bm in [start, end) is set.
func HasBitInRange(bm []uint64, start, end int) bool {
	if start >= end {
		return false
	}
	sw, ew := start/64, (end-1)/64
	if sw == ew {
		return bm[sw]&wordMask(start-sw*64, end-sw*64) != 0
	}
	if bm[sw]&wordMask(start-sw*64, 64) != 0 {
		return true
	}
	if bm[ew]&wordMask(0, end-ew*64) != 0 {
		return true
	}
	return anyNonZero(bm[sw+1 : ew])
}

// PopcountPrefix returns the number of set bits among the first nbits bits of bm.
func PopcountPrefix(bm []uint64, nbits int) int {
	full := nbits / 64
	n := Popcount(bm[:full])
	if rem := nbits % 64; rem > 0 {
		n += bits.OnesCount64(bm[full] & (^uint64(0) >> uint(64-rem)))
	}
	return n
}

// SelectBit returns the position of the rank-th set bit of bm, counting from
// zero, or -1 when bm holds fewer bits than that.
func SelectBit(bm []uint64, rank int) int {
	for i, w := range bm {
		c := bits.OnesCount64(w)
		if rank < c {
			return i*64 + selectInWord(w, rank)
		}
		rank -= c
	}
	return -1
}

// selectInWord returns the position of the rank-th set bit within one word.
func selectInWord(w uint64, rank int) int {
	for ; rank > 0; rank-- {
		w &= w - 1
	}
	return bits.TrailingZeros64(w)
}
