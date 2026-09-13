//go:build goexperiment.simd && amd64

// Package simdops holds the vector primitives shared by every Roaring bitmap
// version in this directory.
//
// Keeping them in one place is deliberate: v1, v2 and v3 differ only in their
// container model, so benchmarks comparing the three measure that model rather
// than accidental differences in the vector code.
//
// The primitives are built on simd/archsimd rather than the portable simd
// package. The portable package offers no population count and no way to turn a
// comparison mask into bits, and nearly every hot Roaring operation needs one or
// the other, so the portable API would only cover the handful of pure bitwise
// loops while the rest dropped to archsimd anyway.
//
// Requires AVX-512 F/BW/VL plus VPOPCNTDQ (population count) and VBMI2
// (compress). Init panics when the CPU lacks them.
package simdops

import (
	"fmt"
	"math/bits"

	"simd/archsimd"
)

// BitmapWords is the number of 64-bit words in a bitmap container, covering the
// full 65536-value key space of a single container.
const BitmapWords = 1024

const (
	words64 = 8  // uint64 lanes in a 512-bit vector
	words16 = 32 // uint16 lanes in a 512-bit vector
)

// lanes16 is the lane index vector [0 1 2 ... 31], used to turn a bit pattern
// into the positions of its set bits.
var lanes16 = func() (v [words16]uint16) {
	for i := range v {
		v[i] = uint16(i)
	}
	return v
}()

func init() {
	missing := make([]string, 0, 3)
	if !archsimd.X86.AVX512() {
		missing = append(missing, "AVX512")
	}
	if !archsimd.X86.AVX512VPOPCNTDQ() {
		missing = append(missing, "AVX512VPOPCNTDQ")
	}
	if !archsimd.X86.AVX512VBMI2() {
		missing = append(missing, "AVX512VBMI2")
	}
	if len(missing) > 0 {
		panic(fmt.Sprintf("simdops: CPU lacks required features: %v", missing))
	}
}

// AndTo stores a&b into dst and returns the number of set bits written.
// dst may alias a or b. All three slices must be the same length.
func AndTo(dst, a, b []uint64) int {
	n := checkTriple(dst, a, b)
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= n; i += words64 {
		v := archsimd.LoadUint64x8(a[i:]).And(archsimd.LoadUint64x8(b[i:]))
		v.Store(dst[i:])
		acc = acc.Add(v.OnesCount())
	}
	total := hsum(acc)
	for ; i < n; i++ {
		dst[i] = a[i] & b[i]
		total += bits.OnesCount64(dst[i])
	}
	return total
}

// OrTo stores a|b into dst and returns the number of set bits written.
func OrTo(dst, a, b []uint64) int {
	n := checkTriple(dst, a, b)
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= n; i += words64 {
		v := archsimd.LoadUint64x8(a[i:]).Or(archsimd.LoadUint64x8(b[i:]))
		v.Store(dst[i:])
		acc = acc.Add(v.OnesCount())
	}
	total := hsum(acc)
	for ; i < n; i++ {
		dst[i] = a[i] | b[i]
		total += bits.OnesCount64(dst[i])
	}
	return total
}

// AndNotTo stores a&^b into dst and returns the number of set bits written.
func AndNotTo(dst, a, b []uint64) int {
	n := checkTriple(dst, a, b)
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= n; i += words64 {
		// archsimd AndNot computes x &^ y, matching the Go operator.
		v := archsimd.LoadUint64x8(a[i:]).AndNot(archsimd.LoadUint64x8(b[i:]))
		v.Store(dst[i:])
		acc = acc.Add(v.OnesCount())
	}
	total := hsum(acc)
	for ; i < n; i++ {
		dst[i] = a[i] &^ b[i]
		total += bits.OnesCount64(dst[i])
	}
	return total
}

// XorTo stores a^b into dst and returns the number of set bits written.
func XorTo(dst, a, b []uint64) int {
	n := checkTriple(dst, a, b)
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= n; i += words64 {
		v := archsimd.LoadUint64x8(a[i:]).Xor(archsimd.LoadUint64x8(b[i:]))
		v.Store(dst[i:])
		acc = acc.Add(v.OnesCount())
	}
	total := hsum(acc)
	for ; i < n; i++ {
		dst[i] = a[i] ^ b[i]
		total += bits.OnesCount64(dst[i])
	}
	return total
}

// AndCardinality returns the number of set bits in a&b without materialising the
// result, for callers that only need the count.
func AndCardinality(a, b []uint64) int {
	n := min(len(a), len(b))
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= n; i += words64 {
		acc = acc.Add(archsimd.LoadUint64x8(a[i:]).And(archsimd.LoadUint64x8(b[i:])).OnesCount())
	}
	total := hsum(acc)
	for ; i < n; i++ {
		total += bits.OnesCount64(a[i] & b[i])
	}
	return total
}

// OrCardinality returns the number of set bits in a|b without materialising it.
func OrCardinality(a, b []uint64) int {
	n := min(len(a), len(b))
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= n; i += words64 {
		acc = acc.Add(archsimd.LoadUint64x8(a[i:]).Or(archsimd.LoadUint64x8(b[i:])).OnesCount())
	}
	total := hsum(acc)
	for ; i < n; i++ {
		total += bits.OnesCount64(a[i] | b[i])
	}
	return total
}

// AndNotCardinality returns the number of set bits in a&^b without materialising it.
func AndNotCardinality(a, b []uint64) int {
	n := min(len(a), len(b))
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= n; i += words64 {
		acc = acc.Add(archsimd.LoadUint64x8(a[i:]).AndNot(archsimd.LoadUint64x8(b[i:])).OnesCount())
	}
	total := hsum(acc)
	for ; i < n; i++ {
		total += bits.OnesCount64(a[i] &^ b[i])
	}
	return total
}

// XorCardinality returns the number of set bits in a^b without materialising it.
func XorCardinality(a, b []uint64) int {
	n := min(len(a), len(b))
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= n; i += words64 {
		acc = acc.Add(archsimd.LoadUint64x8(a[i:]).Xor(archsimd.LoadUint64x8(b[i:])).OnesCount())
	}
	total := hsum(acc)
	for ; i < n; i++ {
		total += bits.OnesCount64(a[i] ^ b[i])
	}
	return total
}

// Popcount returns the number of set bits in a.
func Popcount(a []uint64) int {
	var acc archsimd.Uint64x8
	i := 0
	for ; i+words64 <= len(a); i += words64 {
		acc = acc.Add(archsimd.LoadUint64x8(a[i:]).OnesCount())
	}
	total := hsum(acc)
	for ; i < len(a); i++ {
		total += bits.OnesCount64(a[i])
	}
	return total
}

// Intersects reports whether a&b has any set bit, stopping at the first hit.
func Intersects(a, b []uint64) bool {
	n := min(len(a), len(b))
	i := 0
	for ; i+words64 <= n; i += words64 {
		v := archsimd.LoadUint64x8(a[i:]).And(archsimd.LoadUint64x8(b[i:]))
		if v.NotEqual(archsimd.BroadcastUint64x8(0)).ToBits() != 0 {
			return true
		}
	}
	for ; i < n; i++ {
		if a[i]&b[i] != 0 {
			return true
		}
	}
	return false
}

// BitmapToArray writes the positions of the set bits of bm, in ascending order,
// into dst and returns how many were written. dst must hold Popcount(bm) values.
func BitmapToArray(dst []uint16, bm []uint64) int {
	idx := archsimd.LoadUint16x32Array(&lanes16)
	n := 0
	var buf [words16]uint16
	for wi, w := range bm {
		if w == 0 {
			continue
		}
		base := uint16(wi * 64)
		for half := 0; half < 2; half++ {
			part := uint32(w >> (32 * half))
			if part == 0 {
				continue
			}
			positions := idx.Add(archsimd.BroadcastUint16x32(base + uint16(32*half)))
			positions.Compress(archsimd.Mask16x32FromBits(part)).StoreArray(&buf)
			n += copy(dst[n:], buf[:bits.OnesCount32(part)])
		}
	}
	return n
}

// hsum adds the eight lanes of an accumulator into a single count.
func hsum(v archsimd.Uint64x8) int {
	var out [words64]uint64
	v.StoreArray(&out)
	return int(out[0] + out[1] + out[2] + out[3] + out[4] + out[5] + out[6] + out[7])
}

func checkTriple(dst, a, b []uint64) int {
	if len(a) != len(b) || len(dst) != len(a) {
		panic(fmt.Sprintf("simdops: length mismatch dst=%d a=%d b=%d", len(dst), len(a), len(b)))
	}
	return len(dst)
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
