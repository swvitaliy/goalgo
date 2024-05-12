package min_path

import (
	"cmp"
	"github.com/emirpasic/gods/v2/sets/treeset"
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

type pair[P cmp.Ordered, T any] struct {
	prior P
	val   T
}

func CmpPair[P cmp.Ordered, T any](a, b pair[P, T]) int {
	if a.prior < b.prior {
		return -1
	}
	if a.prior > b.prior {
		return 1
	}
	return 0
}

func makePair[P cmp.Ordered, T any](p P, v T) pair[P, T] {
	return pair[P, T]{p, v}
}

func Dijkstra(a [][]edge, s, t int) (int, []int) {
	n := len(a)

	d := make([]int, n)
	goalSlices.Fill(d, MaxInt)
	d[s] = 0

	p := make([]int, n)
	p[s] = -1

	q := treeset.NewWith[pair[int, int]](CmpPair[int, int])
	q.Add(makePair(0, s))
	for q.Size() > 0 {
		it := q.Iterator()
		it.First()
		w, v := it.Value().prior, it.Value().val
		q.Remove(it.Value())

		if d[v] < w {
			continue
		}
		for _, e := range a[v] {
			if d[v]+e.weight < d[e.to] {
				q.Remove(makePair(d[e.to], e.to))
				d[e.to] = d[v] + e.weight
				p[e.to] = v
				q.Add(makePair(d[e.to], e.to))
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
