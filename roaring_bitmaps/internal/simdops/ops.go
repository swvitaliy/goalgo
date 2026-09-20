package simdops

// Ops exposes the package functions as methods, so a caller can take them
// through an interface and swap them out. The type is empty and every method
// forwards to the package-level function of the same name.
type Ops struct{}

// AndTo forwards to [AndTo].
func (Ops) AndTo(dst, a, b []uint64) int { return AndTo(dst, a, b) }

// OrTo forwards to [OrTo].
func (Ops) OrTo(dst, a, b []uint64) int { return OrTo(dst, a, b) }

// AndNotTo forwards to [AndNotTo].
func (Ops) AndNotTo(dst, a, b []uint64) int { return AndNotTo(dst, a, b) }

// XorTo forwards to [XorTo].
func (Ops) XorTo(dst, a, b []uint64) int { return XorTo(dst, a, b) }

// AndCardinality forwards to [AndCardinality].
func (Ops) AndCardinality(a, b []uint64) int { return AndCardinality(a, b) }

// Popcount forwards to [Popcount].
func (Ops) Popcount(a []uint64) int { return Popcount(a) }

// PopcountPrefix forwards to [PopcountPrefix].
func (Ops) PopcountPrefix(bm []uint64, nbits int) int { return PopcountPrefix(bm, nbits) }

// SelectBit forwards to [SelectBit].
func (Ops) SelectBit(bm []uint64, rank int) int { return SelectBit(bm, rank) }

// Intersects forwards to [Intersects].
func (Ops) Intersects(a, b []uint64) bool { return Intersects(a, b) }

// HasBitInRange forwards to [HasBitInRange].
func (Ops) HasBitInRange(bm []uint64, start, end int) bool { return HasBitInRange(bm, start, end) }

// BitmapToArray forwards to [BitmapToArray].
func (Ops) BitmapToArray(dst []uint16, bm []uint64) int { return BitmapToArray(dst, bm) }

// IntersectArrays forwards to [IntersectArrays].
func (Ops) IntersectArrays(dst, a, b []uint16) int { return IntersectArrays(dst, a, b) }

// UnionArrays forwards to [UnionArrays].
func (Ops) UnionArrays(dst, a, b []uint16) int { return UnionArrays(dst, a, b) }

// DifferenceArrays forwards to [DifferenceArrays].
func (Ops) DifferenceArrays(dst, a, b []uint16) int { return DifferenceArrays(dst, a, b) }

// XorArrays forwards to [XorArrays].
func (Ops) XorArrays(dst, a, b []uint16) int { return XorArrays(dst, a, b) }

// SetRange forwards to [SetRange].
func (Ops) SetRange(bm []uint64, start, end int) { SetRange(bm, start, end) }

// ClearRange forwards to [ClearRange].
func (Ops) ClearRange(bm []uint64, start, end int) { ClearRange(bm, start, end) }

// FlipRange forwards to [FlipRange].
func (Ops) FlipRange(bm []uint64, start, end int) { FlipRange(bm, start, end) }

// CopyRange forwards to [CopyRange].
func (Ops) CopyRange(dst, src []uint64, start, end int) { CopyRange(dst, src, start, end) }
