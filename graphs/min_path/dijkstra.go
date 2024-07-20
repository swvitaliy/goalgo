package min_path

import (
	"github.com/emirpasic/gods/v2/sets/treeset"
	"goalgo/common"
	goalSlices "goalgo/slices"
)

const (
	MaxUInt = ^uint(0)
	MaxInt  = int(MaxUInt >> 1)
	MinInt  = -MaxInt - 1
)

type edge struct {
	to     int
	weight int
}

func Dijkstra(a [][]edge, s, t int) (int, []int) {
	n := len(a)

	d := make([]int, n)
	goalSlices.Fill(d, MaxInt)
	d[s] = 0

	p := make([]int, n)
	p[s] = -1

	q := treeset.NewWith(common.PriorPairCmp[int, int])
	q.Add(common.MakePriorPair(0, s))
	for q.Size() > 0 {
		it := q.Iterator()
		it.First()
		w, v := it.Value().Prior, it.Value().Val
		q.Remove(it.Value())

		if d[v] < w {
			continue
		}
		for _, e := range a[v] {
			if d[v]+e.weight < d[e.to] {
				q.Remove(common.MakePriorPair(d[e.to], e.to))
				d[e.to] = d[v] + e.weight
				p[e.to] = v
				q.Add(common.MakePriorPair(d[e.to], e.to))
			}
		}
	}

	ans := make([]int, 0)
	j := t
	for p[j] != -1 {
		ans = append(ans, j)
		j = p[j]
	}

	goalSlices.Reverse(ans)
	return d[t], ans
}

func DijkstraBidirectional(a [][]edge, s, t int) {

}
