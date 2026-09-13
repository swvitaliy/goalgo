//go:build goexperiment.simd && amd64

package roaringv3

import (
	"slices"

	"goalgo/roaring_bitmaps/internal/simdops"
)

// Rank returns how many values in the set are less than or equal to v.
func (b *Bitmap) Rank(v uint32) int {
	n := 0
	key := hi(v)
	for i, k := range b.keys {
		if k > key {
			break
		}
		if k < key {
			n += int(b.conts[i].card)
			continue
		}
		n += b.conts[i].rank(lo(v))
		break
	}
	return n
}

// Select returns the i-th smallest value in the set, counting from zero, and
// reports whether the set holds that many values.
func (b *Bitmap) Select(i int) (uint32, bool) {
	if i < 0 {
		return 0, false
	}
	for idx, c := range b.conts {
		if i < int(c.card) {
			return uint32(b.keys[idx])<<16 | uint32(c.selectAt(i)), true
		}
		i -= int(c.card)
	}
	return 0, false
}

// Minimum returns the smallest value in the set.
func (b *Bitmap) Minimum() (uint32, bool) { return b.Select(0) }

// Maximum returns the largest value in the set.
func (b *Bitmap) Maximum() (uint32, bool) {
	if len(b.conts) == 0 {
		return 0, false
	}
	last := len(b.conts) - 1
	c := b.conts[last]
	return uint32(b.keys[last])<<16 | uint32(c.selectAt(int(c.card)-1)), true
}

// rank returns how many values of the container are less than or equal to v.
func (c *container) rank(v uint16) int {
	switch c.kind {
	case kindBitmap:
		return simdops.PopcountPrefix(c.bm[:], int(v)+1)

	case kindRun:
		n := 0
		for _, iv := range c.runs {
			if iv.start > v {
				break
			}
			if iv.last <= v {
				n += iv.cardinality()
				continue
			}
			n += int(v) - int(iv.start) + 1
			break
		}
		return n

	default:
		i, found := slices.BinarySearch(c.arr, v)
		if found {
			return i + 1
		}
		return i
	}
}

// selectAt returns the i-th smallest value of the container, counting from zero.
// i must be less than the container's cardinality.
func (c *container) selectAt(i int) uint16 {
	switch c.kind {
	case kindBitmap:
		return uint16(simdops.SelectBit(c.bm[:], i))

	case kindRun:
		for _, iv := range c.runs {
			n := iv.cardinality()
			if i < n {
				return uint16(int(iv.start) + i)
			}
			i -= n
		}
		panic("roaring: selectAt past the end of a run container")

	default:
		return c.arr[i]
	}
}
