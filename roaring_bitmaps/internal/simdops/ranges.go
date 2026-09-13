//go:build goexperiment.simd && amd64

package simdops

import "simd/archsimd"

// The range helpers all take a half-open bit range [start, end). Whole words in
// the middle are the part worth vectorising; the two partial words at the edges
// are handled with masks.

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

func fillWords(dst []uint64, v uint64) {
	vec := archsimd.BroadcastUint64x8(v)
	i := 0
	for ; i+words64 <= len(dst); i += words64 {
		vec.Store(dst[i:])
	}
	for ; i < len(dst); i++ {
		dst[i] = v
	}
}

func flipWords(dst []uint64) {
	ones := archsimd.BroadcastUint64x8(^uint64(0))
	i := 0
	for ; i+words64 <= len(dst); i += words64 {
		archsimd.LoadUint64x8(dst[i:]).Xor(ones).Store(dst[i:])
	}
	for ; i < len(dst); i++ {
		dst[i] = ^dst[i]
	}
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

func anyNonZero(words []uint64) bool {
	zero := archsimd.BroadcastUint64x8(0)
	i := 0
	for ; i+words64 <= len(words); i += words64 {
		if archsimd.LoadUint64x8(words[i:]).NotEqual(zero).ToBits() != 0 {
			return true
		}
	}
	for ; i < len(words); i++ {
		if words[i] != 0 {
			return true
		}
	}
	return false
}
