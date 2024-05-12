package mst

import (
	"goalgo/lca"
	"sort"
)

type edge struct {
	from, to int
	w        wei
}

// KraskalMST for minimum spanning tree
func KraskalMST(g []edge, n int) []edge {
	m := len(g)
	res := make([]edge, 0, m)
	sort.Slice(g, func(i, j int) bool {
		return g[i].w < g[j].w
	})

	d := lca.NewDSU(n)
	for i := 0; i < m; i++ {
		if d.FindSet(g[i].from) != d.FindSet(g[i].to) {
			res = append(res, g[i])
			d.UnionSets(g[i].from, g[i].to, g[i].from)
		}
	}

	return res
}
