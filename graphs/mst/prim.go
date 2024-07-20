package mst

import (
	"github.com/emirpasic/gods/v2/sets/treeset"
	"goalgo/common"
	"goalgo/limits"
	"goalgo/slices"
	"math"
)

type ndx int
type wei int

func init() {
	limits.AddLimits[ndx](math.MinInt, math.MaxInt)
	limits.AddLimits[wei](math.MinInt, math.MaxInt)
}

// adjEdge for adjacent matrix graph representation
type adjEdge struct {
	t ndx
	w wei
}

// PrimMST for minimum spanning tree
func PrimMST(s ndx, g [][]adjEdge) []ndx {
	n := len(g)

	if n == 0 {
		return []ndx{}
	}

	if s >= ndx(n) {
		return []ndx{}
	}

	d := make([]wei, n)
	slices.Fill(d, limits.MaxValue[wei]())
	d[s] = 0

	p := make([]ndx, n)
	p[s] = -1

	q := treeset.NewWith(common.PriorPairCmp[wei, ndx])
	q.Add(common.MakePriorPair(wei(0), s))
	for q.Size() > 0 {
		it := q.Iterator()
		it.First()
		v := it.Value().Val
		q.Remove(it.Value())
		for _, e := range g[v] {
			if e.w < d[e.t] {
				q.Remove(common.MakePriorPair(d[e.t], e.t))
				d[e.t] = e.w
				p[e.t] = v
				q.Add(common.MakePriorPair(d[e.t], e.t))
			}
		}
	}

	return p
}
