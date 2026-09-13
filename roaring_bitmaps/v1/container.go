package roaringv1

import (
	"slices"

	"goalgo/roaring_bitmaps/internal/simdops"
)

// arrayMax is the largest cardinality kept as a sorted array. Above it a bitmap
// costs less: 4097 values as uint16 already exceed the 8KiB a bitmap container
// always occupies.
const arrayMax = 4096

type kind uint8

const (
	kindArray kind = iota
	kindBitmap
)

// bitmapWords is the fixed shape of a bitmap container: a container always covers
// the whole 65536-value key space, so a pointer to the array carries the same
// information as a slice header in a third of the space.
type bitmapWords = [simdops.BitmapWords]uint64

// container holds the low 16 bits of every value sharing one high-16-bit key.
// Both representations live in one struct rather than behind an interface so that
// operations dispatch on a field test instead of an indirect call.
//
// The field order and widths are deliberate. Sparse data produces tens of
// thousands of these per operation, so the struct is kept inside a small size
// class: card is int32 because a container never holds more than 65536 values.
type container struct {
	arr  []uint16     // kindArray, sorted and duplicate-free
	bm   *bitmapWords // kindBitmap
	card int32
	kind kind
}

func newArrayContainer(capacity int) *container {
	return &container{kind: kindArray, arr: make([]uint16, 0, capacity)}
}

func newBitmapContainer() *container {
	return &container{kind: kindBitmap, bm: new(bitmapWords)}
}

func (c *container) clone() *container {
	out := &container{kind: c.kind, card: c.card}
	if c.kind == kindArray {
		out.arr = slices.Clone(c.arr)
	} else {
		out.bm = new(bitmapWords)
		*out.bm = *c.bm
	}
	return out
}

func (c *container) contains(v uint16) bool {
	if c.kind == kindBitmap {
		return c.bm[v/64]&(1<<(v%64)) != 0
	}
	_, found := slices.BinarySearch(c.arr, v)
	return found
}

func (c *container) add(v uint16) bool {
	if c.kind == kindBitmap {
		if !c.setBit(v) {
			return false
		}
		c.card++
		return true
	}

	i, found := slices.BinarySearch(c.arr, v)
	if found {
		return false
	}
	c.arr = slices.Insert(c.arr, i, v)
	c.card++
	if c.card > arrayMax {
		c.convertToBitmap()
	}
	return true
}

func (c *container) remove(v uint16) bool {
	if c.kind == kindBitmap {
		if !c.clearBit(v) {
			return false
		}
		c.card--
		if c.card <= arrayMax {
			c.convertToArray()
		}
		return true
	}

	i, found := slices.BinarySearch(c.arr, v)
	if !found {
		return false
	}
	c.arr = slices.Delete(c.arr, i, i+1)
	c.card--
	return true
}

// setBit, clearBit and flipBit report whether the bit changed to set, so callers
// can keep the cardinality up to date without recounting.
func (c *container) setBit(v uint16) bool {
	w, m := v/64, uint64(1)<<(v%64)
	if c.bm[w]&m != 0 {
		return false
	}
	c.bm[w] |= m
	return true
}

func (c *container) clearBit(v uint16) bool {
	w, m := v/64, uint64(1)<<(v%64)
	if c.bm[w]&m == 0 {
		return false
	}
	c.bm[w] &^= m
	return true
}

func (c *container) flipBit(v uint16) bool {
	w, m := v/64, uint64(1)<<(v%64)
	c.bm[w] ^= m
	return c.bm[w]&m != 0
}

func (c *container) convertToBitmap() {
	bm := new(bitmapWords)
	for _, v := range c.arr {
		bm[v/64] |= 1 << (v % 64)
	}
	c.kind, c.bm, c.arr = kindBitmap, bm, nil
}

func (c *container) convertToArray() {
	arr := make([]uint16, c.card)
	simdops.BitmapToArray(arr, c.bm[:])
	c.kind, c.arr, c.bm = kindArray, arr, nil
}

// normalize picks the cheaper representation for the cardinality a container
// ended up with, and reports nil for an empty result so callers can drop the key.
func (c *container) normalize() *container {
	switch {
	case c == nil || c.card == 0:
		return nil
	case c.kind == kindArray && c.card > arrayMax:
		c.convertToBitmap()
	case c.kind == kindBitmap && c.card <= arrayMax:
		c.convertToArray()
	}
	return c
}

// appendValues appends every value of the container to dst, offset by base.
func (c *container) appendValues(dst []uint32, base uint32) []uint32 {
	if c.kind == kindArray {
		for _, v := range c.arr {
			dst = append(dst, base|uint32(v))
		}
		return dst
	}

	buf := make([]uint16, c.card)
	simdops.BitmapToArray(buf, c.bm[:])
	for _, v := range buf {
		dst = append(dst, base|uint32(v))
	}
	return dst
}

func andContainers(a, b *container) *container {
	switch {
	case a.kind == kindArray && b.kind == kindArray:
		out := newArrayContainer(int(min(a.card, b.card)))
		out.arr = out.arr[:min(a.card, b.card)]
		out.card = int32(simdops.IntersectArrays(out.arr, a.arr, b.arr))
		out.arr = out.arr[:out.card]
		return out.normalize()

	case a.kind == kindArray:
		return filterArray(a, b, true).normalize()

	case b.kind == kindArray:
		return filterArray(b, a, true).normalize()

	default:
		out := newBitmapContainer()
		out.card = int32(simdops.AndTo(out.bm[:], a.bm[:], b.bm[:]))
		return out.normalize()
	}
}

func orContainers(a, b *container) *container {
	switch {
	case a.kind == kindArray && b.kind == kindArray:
		out := newArrayContainer(int(a.card + b.card))
		out.arr = out.arr[:a.card+b.card]
		out.card = int32(simdops.UnionArrays(out.arr, a.arr, b.arr))
		out.arr = out.arr[:out.card]
		return out.normalize()

	case a.kind == kindArray:
		return mergeArrayIntoBitmap(b, a).normalize()

	case b.kind == kindArray:
		return mergeArrayIntoBitmap(a, b).normalize()

	default:
		out := newBitmapContainer()
		out.card = int32(simdops.OrTo(out.bm[:], a.bm[:], b.bm[:]))
		return out.normalize()
	}
}

func andNotContainers(a, b *container) *container {
	switch {
	case a.kind == kindArray && b.kind == kindArray:
		out := newArrayContainer(int(a.card))
		out.arr = out.arr[:a.card]
		out.card = int32(simdops.DifferenceArrays(out.arr, a.arr, b.arr))
		out.arr = out.arr[:out.card]
		return out.normalize()

	case a.kind == kindArray:
		return filterArray(a, b, false).normalize()

	case b.kind == kindArray:
		out := a.clone()
		for _, v := range b.arr {
			if out.clearBit(v) {
				out.card--
			}
		}
		return out.normalize()

	default:
		out := newBitmapContainer()
		out.card = int32(simdops.AndNotTo(out.bm[:], a.bm[:], b.bm[:]))
		return out.normalize()
	}
}

func xorContainers(a, b *container) *container {
	switch {
	case a.kind == kindArray && b.kind == kindArray:
		out := newArrayContainer(int(a.card + b.card))
		out.arr = out.arr[:a.card+b.card]
		out.card = int32(simdops.XorArrays(out.arr, a.arr, b.arr))
		out.arr = out.arr[:out.card]
		return out.normalize()

	case a.kind == kindArray:
		return flipArrayInBitmap(b, a).normalize()

	case b.kind == kindArray:
		return flipArrayInBitmap(a, b).normalize()

	default:
		out := newBitmapContainer()
		out.card = int32(simdops.XorTo(out.bm[:], a.bm[:], b.bm[:]))
		return out.normalize()
	}
}

func intersects(a, b *container) bool {
	switch {
	case a.kind == kindBitmap && b.kind == kindBitmap:
		return simdops.Intersects(a.bm[:], b.bm[:])

	case a.kind == kindArray && b.kind == kindBitmap:
		return anyContained(a.arr, b)

	case a.kind == kindBitmap && b.kind == kindArray:
		return anyContained(b.arr, a)

	default:
		small, large := a, b
		if small.card > large.card {
			small, large = large, small
		}
		for _, v := range small.arr {
			if large.contains(v) {
				return true
			}
		}
		return false
	}
}

func anyContained(values []uint16, bm *container) bool {
	for _, v := range values {
		if bm.bm[v/64]&(1<<(v%64)) != 0 {
			return true
		}
	}
	return false
}

// filterArray keeps the values of the array container arr whose membership in the
// bitmap container bm equals want, producing an array container.
func filterArray(arr, bm *container, want bool) *container {
	out := newArrayContainer(int(arr.card))
	for _, v := range arr.arr {
		if (bm.bm[v/64]&(1<<(v%64)) != 0) == want {
			out.arr = append(out.arr, v)
		}
	}
	out.card = int32(len(out.arr))
	return out
}

func mergeArrayIntoBitmap(bm, arr *container) *container {
	out := bm.clone()
	for _, v := range arr.arr {
		if out.setBit(v) {
			out.card++
		}
	}
	return out
}

func flipArrayInBitmap(bm, arr *container) *container {
	out := bm.clone()
	for _, v := range arr.arr {
		if out.flipBit(v) {
			out.card++
		} else {
			out.card--
		}
	}
	return out
}
