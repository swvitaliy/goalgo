//go:build goexperiment.simd && amd64

package simdops

import (
	"fmt"
	"math/bits"

	"simd/archsimd"
)

// Vectorized reports that this build uses the AVX-512 primitives.
const Vectorized = true

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

	return n + intersectLinear(dst[n:], a[i:], b[j:])
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
