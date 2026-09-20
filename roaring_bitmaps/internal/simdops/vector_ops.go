//go:build goexperiment.simd && amd64

package simdops

// VectorOps exposes the AVX-512 primitives as methods, so a caller can take
// them through an interface and swap them out. It exists only in the vector
// build; the type is empty and every method forwards to the package-level
// function of the same name, which in this build is the vectorised one for the
// hot primitives and the shared word-level code for the rest.
type VectorOps struct{}

// AndTo forwards to [AndTo].
func (VectorOps) AndTo(dst, a, b []uint64) int { return AndTo(dst, a, b) }

// OrTo forwards to [OrTo].
func (VectorOps) OrTo(dst, a, b []uint64) int { return OrTo(dst, a, b) }

// AndNotTo forwards to [AndNotTo].
func (VectorOps) AndNotTo(dst, a, b []uint64) int { return AndNotTo(dst, a, b) }

// XorTo forwards to [XorTo].
func (VectorOps) XorTo(dst, a, b []uint64) int { return XorTo(dst, a, b) }

// AndCardinality forwards to [AndCardinality].
func (VectorOps) AndCardinality(a, b []uint64) int { return AndCardinality(a, b) }

// Popcount forwards to [Popcount].
func (VectorOps) Popcount(a []uint64) int { return Popcount(a) }

// PopcountPrefix forwards to [PopcountPrefix].
func (VectorOps) PopcountPrefix(bm []uint64, nbits int) int { return PopcountPrefix(bm, nbits) }

// SelectBit forwards to [SelectBit].
func (VectorOps) SelectBit(bm []uint64, rank int) int { return SelectBit(bm, rank) }

// Intersects forwards to [Intersects].
func (VectorOps) Intersects(a, b []uint64) bool { return Intersects(a, b) }

// HasBitInRange forwards to [HasBitInRange].
func (VectorOps) HasBitInRange(bm []uint64, start, end int) bool {
	return HasBitInRange(bm, start, end)
}

// BitmapToArray forwards to [BitmapToArray].
func (VectorOps) BitmapToArray(dst []uint16, bm []uint64) int { return BitmapToArray(dst, bm) }

// IntersectArrays forwards to [IntersectArrays].
func (VectorOps) IntersectArrays(dst, a, b []uint16) int { return IntersectArrays(dst, a, b) }

// UnionArrays forwards to [UnionArrays].
func (VectorOps) UnionArrays(dst, a, b []uint16) int { return UnionArrays(dst, a, b) }

// DifferenceArrays forwards to [DifferenceArrays].
func (VectorOps) DifferenceArrays(dst, a, b []uint16) int { return DifferenceArrays(dst, a, b) }

// XorArrays forwards to [XorArrays].
func (VectorOps) XorArrays(dst, a, b []uint16) int { return XorArrays(dst, a, b) }

// SetRange forwards to [SetRange].
func (VectorOps) SetRange(bm []uint64, start, end int) { SetRange(bm, start, end) }

// ClearRange forwards to [ClearRange].
func (VectorOps) ClearRange(bm []uint64, start, end int) { ClearRange(bm, start, end) }

// FlipRange forwards to [FlipRange].
func (VectorOps) FlipRange(bm []uint64, start, end int) { FlipRange(bm, start, end) }

// CopyRange forwards to [CopyRange].
func (VectorOps) CopyRange(dst, src []uint64, start, end int) { CopyRange(dst, src, start, end) }
