package min_path

import (
	"goalgo/pq"
	"goalgo/slices"
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
	for i := 0; i < n; i++ {
		d[i] = MaxInt
	}
	d[s] = 0
	p := make([]int, n)
	p[s] = -1
	q := pq.NewPQ[int]()
	q.Enqueue(0, s)
	for len(q) > 0 {
		v, l := q.Dequeue()
		if d[v] < l {
			continue
		}
		for _, e := range a[v] {
			if d[v]+e.weight < d[e.to] {
				d[e.to] = d[v] + e.weight
				q.Enqueue(d[e.to], e.to)
				p[e.to] = v
			}
		}
	}

	ans := make([]int, 0)
	j := t
	for p[j] != -1 {
		ans = append(ans, j)
		j = p[j]
	}

	slices.Reverse(ans)
	return d[t], ans
}

func DijkstraBidirectional(a [][]edge, s, t int) {

}
