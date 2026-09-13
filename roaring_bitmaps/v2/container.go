//go:build goexperiment.simd && amd64

package roaringv2

import (
	"slices"

	"goalgo/roaring_bitmaps/internal/simdops"
)

const (
	// arrayMax is the largest cardinality kept as a sorted array: above it the
	// fixed 8KiB of a bitmap costs less than 2 bytes per value.
	arrayMax = 4096
	// bitmapBytes is what a bitmap container always occupies.
	bitmapBytes = simdops.BitmapWords * 8
)

type kind uint8

const (
	kindArray kind = iota
	kindBitmap
	kindRun
)

// bitmapWords is the fixed shape of a bitmap container: a container always covers
// the whole 65536-value key space, so a pointer to the array carries the same
// information as a slice header in a third of the space.
type bitmapWords = [simdops.BitmapWords]uint64

// container holds the low 16 bits of every value sharing one high-16-bit key, in
// whichever of the three encodings is cheapest for its contents.
//
// The field order and widths are deliberate. Sparse data produces tens of
// thousands of these per operation, so the struct is kept inside a small size
// class: card is int32 because a container never holds more than 65536 values,
// and the bitmap is a pointer rather than a slice. Carrying the third encoding
// costs one extra slice header over the two-encoding container of v1.
type container struct {
	arr  []uint16     // kindArray, sorted and duplicate-free
	runs []interval   // kindRun, sorted, disjoint and non-adjacent
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

func newRunContainer(runs []interval) *container {
	return &container{kind: kindRun, card: int32(runsCardinality(runs)), runs: runs}
}

func (c *container) clone() *container {
	out := &container{kind: c.kind, card: c.card}
	switch c.kind {
	case kindArray:
		out.arr = slices.Clone(c.arr)
	case kindBitmap:
		out.bm = new(bitmapWords)
		*out.bm = *c.bm
	case kindRun:
		out.runs = slices.Clone(c.runs)
	}
	return out
}

func (c *container) contains(v uint16) bool {
	switch c.kind {
	case kindBitmap:
		return c.bm[v/64]&(1<<(v%64)) != 0
	case kindRun:
		return runContains(c.runs, v)
	default:
		_, found := slices.BinarySearch(c.arr, v)
		return found
	}
}

func (c *container) add(v uint16) bool {
	switch c.kind {
	case kindBitmap:
		if !c.setBit(v) {
			return false
		}
		c.card++

	case kindRun:
		runs, added := runAdd(c.runs, v)
		if !added {
			return false
		}
		c.runs = runs
		c.card++
		c.normalize()

	default:
		i, found := slices.BinarySearch(c.arr, v)
		if found {
			return false
		}
		c.arr = slices.Insert(c.arr, i, v)
		c.card++
		if c.card > arrayMax {
			c.convertToBitmap()
		}
	}
	return true
}

func (c *container) remove(v uint16) bool {
	switch c.kind {
	case kindBitmap:
		if !c.clearBit(v) {
			return false
		}
		c.card--
		if c.card <= arrayMax {
			c.convertToArray()
		}

	case kindRun:
		runs, removed := runRemove(c.runs, v)
		if !removed {
			return false
		}
		c.runs = runs
		c.card--
		c.normalize()

	default:
		i, found := slices.BinarySearch(c.arr, v)
		if !found {
			return false
		}
		c.arr = slices.Delete(c.arr, i, i+1)
		c.card--
	}
	return true
}

// setBit, clearBit and flipBit report whether the bit ended up set, so callers
// can track cardinality without recounting.
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

func (c *container) convertToBitmap() {
	switch c.kind {
	case kindArray:
		bm := new(bitmapWords)
		for _, v := range c.arr {
			bm[v/64] |= 1 << (v % 64)
		}
		c.bm, c.arr = bm, nil
	case kindRun:
		c.bm, c.runs = runsToBitmap(c.runs), nil
	}
	c.kind = kindBitmap
}

func (c *container) convertToArray() {
	arr := make([]uint16, c.card)
	switch c.kind {
	case kindBitmap:
		simdops.BitmapToArray(arr, c.bm[:])
		c.bm = nil
	case kindRun:
		runsToArray(arr, c.runs)
		c.runs = nil
	}
	c.kind, c.arr = kindArray, arr
}

// runOptimize switches the container to run encoding when that encoding is the
// smallest of the three. It is only ever called explicitly, so array and bitmap
// containers keep their shape until the caller asks.
func (c *container) runOptimize() {
	if c.kind == kindRun {
		return
	}

	var nruns int
	if c.kind == kindArray {
		nruns = countRunsArray(c.arr)
	} else {
		nruns = countRunsBitmap(c.bm[:])
	}
	if runsBytes(nruns) >= c.currentBytes() {
		return
	}

	if c.kind == kindArray {
		c.runs = arrayToRuns(c.arr)
		c.arr = nil
	} else {
		c.runs = bitmapToRuns(c.bm[:])
		c.bm = nil
	}
	c.kind = kindRun
}

func (c *container) currentBytes() int {
	switch c.kind {
	case kindArray:
		return 2 * int(c.card)
	case kindBitmap:
		return bitmapBytes
	default:
		return runsBytes(len(c.runs))
	}
}

// normalize picks the cheapest representation for the contents a container ended
// up with, and reports nil for an empty result so callers can drop the key.
//
// A run container is kept only while runs remain the smallest encoding; array and
// bitmap containers convert between each other at arrayMax but never turn into
// runs on their own.
func (c *container) normalize() *container {
	if c == nil || c.card == 0 {
		return nil
	}

	switch c.kind {
	case kindRun:
		if runsBytes(len(c.runs)) > min(2*int(c.card), bitmapBytes) {
			if c.card <= arrayMax {
				c.convertToArray()
			} else {
				c.convertToBitmap()
			}
		}
	case kindArray:
		if c.card > arrayMax {
			c.convertToBitmap()
		}
	case kindBitmap:
		if c.card <= arrayMax {
			c.convertToArray()
		}
	}
	return c
}

// appendValues appends every value of the container to dst, offset by base.
func (c *container) appendValues(dst []uint32, base uint32) []uint32 {
	switch c.kind {
	case kindArray:
		for _, v := range c.arr {
			dst = append(dst, base|uint32(v))
		}
	case kindRun:
		for _, iv := range c.runs {
			for v := int(iv.start); v <= int(iv.last); v++ {
				dst = append(dst, base|uint32(v))
			}
		}
	default:
		buf := make([]uint16, c.card)
		simdops.BitmapToArray(buf, c.bm[:])
		for _, v := range buf {
			dst = append(dst, base|uint32(v))
		}
	}
	return dst
}
