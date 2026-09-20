package roaring_bitmaps

import "slices"

// The binary operations dispatch on the encodings of both operands. Run operands
// are handled by interval arithmetic where both sides are runs, and otherwise by
// whichever side already has the cheaper structure to walk: array values are
// probed against runs, and bitmaps are edited range by range.

func andContainers(ops Ops, a, b *container) *container {
	switch {
	case a.kind == kindRun && b.kind == kindRun:
		return newRunContainer(intersectRuns(a.runs, b.runs)).normalize(ops)

	case a.kind == kindRun && b.kind == kindArray:
		return andRunArray(a, b).normalize(ops)
	case a.kind == kindArray && b.kind == kindRun:
		return andRunArray(b, a).normalize(ops)

	case a.kind == kindRun && b.kind == kindBitmap:
		return andRunBitmap(ops, a, b).normalize(ops)
	case a.kind == kindBitmap && b.kind == kindRun:
		return andRunBitmap(ops, b, a).normalize(ops)

	case a.kind == kindArray && b.kind == kindArray:
		out := newArrayContainer(int(min(a.card, b.card)))
		out.arr = out.arr[:min(a.card, b.card)]
		out.card = int32(ops.IntersectArrays(out.arr, a.arr, b.arr))
		out.arr = out.arr[:out.card]
		return out.normalize(ops)

	case a.kind == kindArray:
		return filterArray(a, b, true).normalize(ops)
	case b.kind == kindArray:
		return filterArray(b, a, true).normalize(ops)

	default:
		out := newBitmapContainer()
		out.card = int32(ops.AndTo(out.bm[:], a.bm[:], b.bm[:]))
		return out.normalize(ops)
	}
}

func orContainers(ops Ops, a, b *container) *container {
	switch {
	case a.kind == kindRun && b.kind == kindRun:
		return newRunContainer(unionRuns(a.runs, b.runs)).normalize(ops)

	case a.kind == kindRun && b.kind == kindArray:
		return newRunContainer(unionRuns(a.runs, arrayToRuns(b.arr))).normalize(ops)
	case a.kind == kindArray && b.kind == kindRun:
		return newRunContainer(unionRuns(b.runs, arrayToRuns(a.arr))).normalize(ops)

	case a.kind == kindRun && b.kind == kindBitmap:
		return orRunBitmap(ops, a, b).normalize(ops)
	case a.kind == kindBitmap && b.kind == kindRun:
		return orRunBitmap(ops, b, a).normalize(ops)

	case a.kind == kindArray && b.kind == kindArray:
		out := newArrayContainer(int(a.card + b.card))
		out.arr = out.arr[:a.card+b.card]
		out.card = int32(ops.UnionArrays(out.arr, a.arr, b.arr))
		out.arr = out.arr[:out.card]
		return out.normalize(ops)

	case a.kind == kindArray:
		return mergeArrayIntoBitmap(b, a).normalize(ops)
	case b.kind == kindArray:
		return mergeArrayIntoBitmap(a, b).normalize(ops)

	default:
		out := newBitmapContainer()
		out.card = int32(ops.OrTo(out.bm[:], a.bm[:], b.bm[:]))
		return out.normalize(ops)
	}
}

func andNotContainers(ops Ops, a, b *container) *container {
	switch {
	case a.kind == kindRun && b.kind == kindRun:
		return newRunContainer(differenceRuns(a.runs, b.runs)).normalize(ops)

	case a.kind == kindRun && b.kind == kindArray:
		return newRunContainer(differenceRuns(a.runs, arrayToRuns(b.arr))).normalize(ops)
	case a.kind == kindArray && b.kind == kindRun:
		return filterArrayByRuns(a, b.runs, false).normalize(ops)

	case a.kind == kindRun && b.kind == kindBitmap:
		out := &container{kind: kindBitmap, bm: runsToBitmap(ops, a.runs)}
		out.card = int32(ops.AndNotTo(out.bm[:], out.bm[:], b.bm[:]))
		return out.normalize(ops)
	case a.kind == kindBitmap && b.kind == kindRun:
		out := a.clone()
		for _, iv := range b.runs {
			ops.ClearRange(out.bm[:], int(iv.start), int(iv.last)+1)
		}
		out.card = int32(ops.Popcount(out.bm[:]))
		return out.normalize(ops)

	case a.kind == kindArray && b.kind == kindArray:
		out := newArrayContainer(int(a.card))
		out.arr = out.arr[:a.card]
		out.card = int32(ops.DifferenceArrays(out.arr, a.arr, b.arr))
		out.arr = out.arr[:out.card]
		return out.normalize(ops)

	case a.kind == kindArray:
		return filterArray(a, b, false).normalize(ops)

	case b.kind == kindArray:
		out := a.clone()
		for _, v := range b.arr {
			if out.clearBit(v) {
				out.card--
			}
		}
		return out.normalize(ops)

	default:
		out := newBitmapContainer()
		out.card = int32(ops.AndNotTo(out.bm[:], a.bm[:], b.bm[:]))
		return out.normalize(ops)
	}
}

func xorContainers(ops Ops, a, b *container) *container {
	switch {
	case a.kind == kindRun && b.kind == kindRun:
		return newRunContainer(xorRuns(a.runs, b.runs)).normalize(ops)

	case a.kind == kindRun && b.kind == kindArray:
		return newRunContainer(xorRuns(a.runs, arrayToRuns(b.arr))).normalize(ops)
	case a.kind == kindArray && b.kind == kindRun:
		return newRunContainer(xorRuns(b.runs, arrayToRuns(a.arr))).normalize(ops)

	case a.kind == kindRun && b.kind == kindBitmap:
		return xorRunBitmap(ops, a, b).normalize(ops)
	case a.kind == kindBitmap && b.kind == kindRun:
		return xorRunBitmap(ops, b, a).normalize(ops)

	case a.kind == kindArray && b.kind == kindArray:
		out := newArrayContainer(int(a.card + b.card))
		out.arr = out.arr[:a.card+b.card]
		out.card = int32(ops.XorArrays(out.arr, a.arr, b.arr))
		out.arr = out.arr[:out.card]
		return out.normalize(ops)

	case a.kind == kindArray:
		return flipArrayInBitmap(b, a).normalize(ops)
	case b.kind == kindArray:
		return flipArrayInBitmap(a, b).normalize(ops)

	default:
		out := newBitmapContainer()
		out.card = int32(ops.XorTo(out.bm[:], a.bm[:], b.bm[:]))
		return out.normalize(ops)
	}
}

func intersects(ops Ops, a, b *container) bool {
	switch {
	case a.kind == kindRun && b.kind == kindRun:
		return runsIntersect(a.runs, b.runs)

	case a.kind == kindRun:
		return runsIntersectContainer(ops, a.runs, b)
	case b.kind == kindRun:
		return runsIntersectContainer(ops, b.runs, a)

	case a.kind == kindBitmap && b.kind == kindBitmap:
		return ops.Intersects(a.bm[:], b.bm[:])

	case a.kind == kindArray && b.kind == kindBitmap:
		return anyContained(a.arr, b)
	case a.kind == kindBitmap && b.kind == kindArray:
		return anyContained(b.arr, a)

	default:
		small, large := a, b
		if small.card > large.card {
			small, large = large, small
		}
		return slices.ContainsFunc(small.arr, large.contains)
	}
}

func runsIntersectContainer(ops Ops, runs []interval, c *container) bool {
	if c.kind == kindArray {
		j := 0
		for _, v := range c.arr {
			for j < len(runs) && runs[j].last < v {
				j++
			}
			if j == len(runs) {
				return false
			}
			if runs[j].start <= v {
				return true
			}
		}
		return false
	}
	for _, iv := range runs {
		if ops.HasBitInRange(c.bm[:], int(iv.start), int(iv.last)+1) {
			return true
		}
	}
	return false
}

// andRunArray keeps the array values covered by a run, walking both in order.
func andRunArray(r, arr *container) *container {
	return filterArrayByRuns(arr, r.runs, true)
}

// filterArrayByRuns keeps the values of an array container whose membership in
// runs equals want.
func filterArrayByRuns(arr *container, runs []interval, want bool) *container {
	out := newArrayContainer(int(arr.card))
	j := 0
	for _, v := range arr.arr {
		for j < len(runs) && runs[j].last < v {
			j++
		}
		covered := j < len(runs) && runs[j].start <= v
		if covered == want {
			out.arr = append(out.arr, v)
		}
	}
	out.card = int32(len(out.arr))
	return out
}

// andRunBitmap copies only the bitmap ranges the runs cover, so untouched words
// stay zero without a full AND against a materialised run bitmap.
func andRunBitmap(ops Ops, r, bmc *container) *container {
	out := newBitmapContainer()
	for _, iv := range r.runs {
		ops.CopyRange(out.bm[:], bmc.bm[:], int(iv.start), int(iv.last)+1)
	}
	out.card = int32(ops.Popcount(out.bm[:]))
	return out
}

func orRunBitmap(ops Ops, r, bmc *container) *container {
	out := bmc.clone()
	for _, iv := range r.runs {
		ops.SetRange(out.bm[:], int(iv.start), int(iv.last)+1)
	}
	out.card = int32(ops.Popcount(out.bm[:]))
	return out
}

func xorRunBitmap(ops Ops, r, bmc *container) *container {
	out := bmc.clone()
	for _, iv := range r.runs {
		ops.FlipRange(out.bm[:], int(iv.start), int(iv.last)+1)
	}
	out.card = int32(ops.Popcount(out.bm[:]))
	return out
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
// bitmap container bm equals want.
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
		w, m := v/64, uint64(1)<<(v%64)
		out.bm[w] ^= m
		if out.bm[w]&m != 0 {
			out.card++
		} else {
			out.card--
		}
	}
	return out
}
