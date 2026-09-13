package roaringv2

import (
	"math/bits"
	"slices"

	"goalgo/roaring_bitmaps/internal/simdops"
)

// interval is one run of consecutive values, inclusive at both ends. Runs inside
// a container are sorted, non-overlapping and never adjacent: two runs that touch
// are always merged into one, which is what keeps the encoding minimal.
type interval struct {
	start uint16
	last  uint16
}

func (iv interval) cardinality() int { return int(iv.last) - int(iv.start) + 1 }

func runsCardinality(runs []interval) int {
	n := 0
	for _, iv := range runs {
		n += iv.cardinality()
	}
	return n
}

// runsBytes is the serialised size of a run container: a count plus two uint16
// per run. It is what decides whether run encoding is worth choosing.
func runsBytes(nruns int) int { return 2 + 4*nruns }

func runContains(runs []interval, v uint16) bool {
	_, found := slices.BinarySearchFunc(runs, v, func(iv interval, v uint16) int {
		switch {
		case iv.last < v:
			return -1
		case iv.start > v:
			return 1
		default:
			return 0
		}
	})
	return found
}

// runAdd inserts v, merging with neighbouring runs where it closes a gap. It
// reports whether the value was absent.
func runAdd(runs []interval, v uint16) ([]interval, bool) {
	i, found := slices.BinarySearchFunc(runs, v, func(iv interval, v uint16) int {
		switch {
		case iv.last < v:
			return -1
		case iv.start > v:
			return 1
		default:
			return 0
		}
	})
	if found {
		return runs, false
	}

	extendsPrev := i > 0 && int(runs[i-1].last)+1 == int(v)
	extendsNext := i < len(runs) && int(runs[i].start)-1 == int(v)

	switch {
	case extendsPrev && extendsNext:
		runs[i-1].last = runs[i].last
		runs = slices.Delete(runs, i, i+1)
	case extendsPrev:
		runs[i-1].last = v
	case extendsNext:
		runs[i].start = v
	default:
		runs = slices.Insert(runs, i, interval{start: v, last: v})
	}
	return runs, true
}

// runRemove deletes v, splitting a run when the value sits in its middle. It
// reports whether the value was present.
func runRemove(runs []interval, v uint16) ([]interval, bool) {
	i, found := slices.BinarySearchFunc(runs, v, func(iv interval, v uint16) int {
		switch {
		case iv.last < v:
			return -1
		case iv.start > v:
			return 1
		default:
			return 0
		}
	})
	if !found {
		return runs, false
	}

	iv := runs[i]
	switch {
	case iv.start == iv.last:
		runs = slices.Delete(runs, i, i+1)
	case iv.start == v:
		runs[i].start = v + 1
	case iv.last == v:
		runs[i].last = v - 1
	default:
		runs[i].last = v - 1
		runs = slices.Insert(runs, i+1, interval{start: v + 1, last: iv.last})
	}
	return runs, true
}

func intersectRuns(a, b []interval) []interval {
	out := make([]interval, 0, min(len(a), len(b)))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		start := max(a[i].start, b[j].start)
		last := min(a[i].last, b[j].last)
		if start <= last {
			out = append(out, interval{start: start, last: last})
		}
		if a[i].last < b[j].last {
			i++
		} else {
			j++
		}
	}
	return out
}

func unionRuns(a, b []interval) []interval {
	out := make([]interval, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		var iv interval
		switch {
		case j == len(b) || (i < len(a) && a[i].start <= b[j].start):
			iv, i = a[i], i+1
		default:
			iv, j = b[j], j+1
		}
		out = appendCoalesced(out, iv)
	}
	return out
}

// appendCoalesced appends iv, merging it into the previous run when the two
// overlap or touch. iv must start no earlier than the previous run.
func appendCoalesced(out []interval, iv interval) []interval {
	if n := len(out); n > 0 && int(out[n-1].last)+1 >= int(iv.start) {
		if iv.last > out[n-1].last {
			out[n-1].last = iv.last
		}
		return out
	}
	return append(out, iv)
}

func differenceRuns(a, b []interval) []interval {
	out := make([]interval, 0, len(a)+len(b))
	j := 0
	for _, r := range a {
		start, last := int(r.start), int(r.last)

		// Runs of b that end before this one begins are of no further use: a is
		// sorted, so later runs start even later.
		for j < len(b) && int(b[j].last) < start {
			j++
		}

		for k := j; k < len(b) && int(b[k].start) <= last; k++ {
			bs, bl := int(b[k].start), int(b[k].last)
			if bs > start {
				out = append(out, interval{start: uint16(start), last: uint16(bs - 1)})
			}
			if bl >= last {
				start = last + 1
				break
			}
			start = bl + 1
		}

		if start <= last {
			out = append(out, interval{start: uint16(start), last: uint16(last)})
		}
	}
	return out
}

func xorRuns(a, b []interval) []interval {
	return unionRuns(differenceRuns(a, b), differenceRuns(b, a))
}

func runsIntersect(a, b []interval) bool {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if max(a[i].start, b[j].start) <= min(a[i].last, b[j].last) {
			return true
		}
		if a[i].last < b[j].last {
			i++
		} else {
			j++
		}
	}
	return false
}

// arrayToRuns collapses a sorted array into runs.
func arrayToRuns(arr []uint16) []interval {
	if len(arr) == 0 {
		return nil
	}

	out := make([]interval, 0, len(arr))
	cur := interval{start: arr[0], last: arr[0]}
	for _, v := range arr[1:] {
		if int(v) == int(cur.last)+1 {
			cur.last = v
			continue
		}
		out = append(out, cur)
		cur = interval{start: v, last: v}
	}
	return append(out, cur)
}

func runsToArray(dst []uint16, runs []interval) int {
	n := 0
	for _, iv := range runs {
		for v := int(iv.start); v <= int(iv.last); v++ {
			dst[n] = uint16(v)
			n++
		}
	}
	return n
}

func runsToBitmap(runs []interval) *bitmapWords {
	bm := new(bitmapWords)
	for _, iv := range runs {
		simdops.SetRange(bm[:], int(iv.start), int(iv.last)+1)
	}
	return bm
}

// bitmapToRuns walks the set bits of bm, emitting one interval per maximal run.
func bitmapToRuns(bm []uint64) []interval {
	out := make([]interval, 0, countRunsBitmap(bm))
	pos := 0
	total := len(bm) * 64
	for pos < total {
		start := nextSetBit(bm, pos)
		if start < 0 {
			break
		}
		end := nextClearBit(bm, start)
		out = append(out, interval{start: uint16(start), last: uint16(end - 1)})
		pos = end
	}
	return out
}

func nextSetBit(bm []uint64, from int) int {
	w := from / 64
	if w >= len(bm) {
		return -1
	}
	if rest := bm[w] >> uint(from%64); rest != 0 {
		return from + bits.TrailingZeros64(rest)
	}
	for w++; w < len(bm); w++ {
		if bm[w] != 0 {
			return w*64 + bits.TrailingZeros64(bm[w])
		}
	}
	return -1
}

// nextClearBit returns the first clear bit at or after from, or len(bm)*64 when
// every remaining bit is set.
func nextClearBit(bm []uint64, from int) int {
	w := from / 64
	if rest := ^bm[w] >> uint(from%64); rest != 0 {
		return from + bits.TrailingZeros64(rest)
	}
	for w++; w < len(bm); w++ {
		if bm[w] != ^uint64(0) {
			return w*64 + bits.TrailingZeros64(^bm[w])
		}
	}
	return len(bm) * 64
}

// countRunsBitmap counts maximal runs of set bits: a run starts at every set bit
// whose predecessor is clear.
func countRunsBitmap(bm []uint64) int {
	runs, prev := 0, uint64(0)
	for _, w := range bm {
		runs += bits.OnesCount64(w &^ (w<<1 | prev>>63))
		prev = w
	}
	return runs
}

// countRunsArray counts maximal runs in a sorted array.
func countRunsArray(arr []uint16) int {
	if len(arr) == 0 {
		return 0
	}
	runs := 1
	for i := 1; i < len(arr); i++ {
		if int(arr[i]) != int(arr[i-1])+1 {
			runs++
		}
	}
	return runs
}
