package roaring_bitmaps

//go:generate mockgen -source=ports.go -destination=ports_mock.go -package=roaring_bitmaps

// Ops is the set of bit-level primitives the bitmap runs on. Every hot loop over
// bitmap words and sorted arrays goes through it, so the implementation can be
// swapped: internal/simdops provides the real one (AVX-512 or plain Go depending
// on the build), and a test can supply a mock to observe or fake the calls.
//
// Slice arguments are borrowed for the duration of the call; the word slices
// always have BitmapWords entries.
type Ops interface {
	// AndTo, OrTo, AndNotTo and XorTo store the word-wise result into dst and
	// return the number of set bits written. dst may alias a or b.
	AndTo(dst, a, b []uint64) int
	OrTo(dst, a, b []uint64) int
	AndNotTo(dst, a, b []uint64) int
	XorTo(dst, a, b []uint64) int

	// AndCardinality returns the number of set bits in a&b without materialising it.
	AndCardinality(a, b []uint64) int
	// Popcount returns the number of set bits in a.
	Popcount(a []uint64) int
	// PopcountPrefix returns the number of set bits among the first nbits bits of bm.
	PopcountPrefix(bm []uint64, nbits int) int
	// SelectBit returns the position of the rank-th set bit of bm, counting from
	// zero, or -1 when bm holds fewer bits than that.
	SelectBit(bm []uint64, rank int) int
	// Intersects reports whether a&b has any set bit.
	Intersects(a, b []uint64) bool
	// HasBitInRange reports whether any bit of bm in [start, end) is set.
	HasBitInRange(bm []uint64, start, end int) bool

	// BitmapToArray writes the positions of the set bits of bm, ascending, into
	// dst and returns how many were written. dst must hold Popcount(bm) values.
	BitmapToArray(dst []uint16, bm []uint64) int

	// IntersectArrays, UnionArrays, DifferenceArrays and XorArrays write the set
	// operation over the sorted, duplicate-free arrays a and b into dst and
	// return how many values were written. dst must not alias a or b and must
	// have room for the largest possible result.
	IntersectArrays(dst, a, b []uint16) int
	UnionArrays(dst, a, b []uint16) int
	DifferenceArrays(dst, a, b []uint16) int
	XorArrays(dst, a, b []uint16) int

	// SetRange, ClearRange and FlipRange edit every bit of bm in [start, end).
	SetRange(bm []uint64, start, end int)
	ClearRange(bm []uint64, start, end int)
	FlipRange(bm []uint64, start, end int)
	// CopyRange copies the bits of src in [start, end) into dst, leaving the
	// bits of dst outside the range untouched.
	CopyRange(dst, src []uint64, start, end int)
}
