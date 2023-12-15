package mst

import (
	"goalgo/limits"
	"goalgo/pq"
	"goalgo/slices"
	"math"
)

type ndx int
type wei int

func init() {
	limits.AddTypeLimits[ndx](limits.Limits{
		MaxValue: math.MaxInt,
		MinValue: math.MinInt,
	})
	limits.AddTypeLimits[wei](limits.Limits{
		MaxValue: math.MaxInt,
		MinValue: math.MinInt,
	})
}

// adjEdge for adjacent matrix graph representation
type adjEdge struct {
	t ndx
	w wei
}

// PrimMST for minimum spanning tree
func PrimMST(s ndx, g [][]adjEdge) []ndx {
	n := len(g)
	d := make([]wei, n)
	slices.Fill(d, limits.MaxValue[wei]())
	d[s] = 0
	p := make([]ndx, n)
	p[s] = -1

	q := pq.NewPriorQueue[ndx, wei]()
	q.Enqueue(0, s)
	for len(q) > 0 {
		v, _ := q.Dequeue()
		for _, e := range g[v] {
			if e.w < d[e.t] {
				// TODO replace with btree
				d[e.t] = e.w
				p[e.t] = v
				q.Enqueue(d[e.t], e.t)
			}
		}
	}

	return p
}
