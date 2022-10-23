package main

import "fmt"

type DSU struct {
	p []int
	r []int
	a []int
	u []bool
}

func NewDSU(n int) *DSU {
	ans := DSU{
		p: make([]int, n),
		r: make([]int, n),
		a: make([]int, n),
		u: make([]bool, n),
	}

	for i := 0; i < n; i++ {
		ans.r[i] = 1
		ans.u[i] = false
	}

	return &ans
}

func (dsu *DSU) FindSet(v int) int {
	if dsu.p[v] == v {
		return dsu.p[v]
	}

	dsu.p[v] = dsu.FindSet(dsu.p[v])
	return dsu.p[v]
}

func (dsu *DSU) UnionSets(u int, v int, x int) {
	u = dsu.FindSet(u)
	v = dsu.FindSet(v)
	if u != v {
		if dsu.r[v] > dsu.r[u] {
			u, v = v, u
		}

		// r[v] >= r[u]
		dsu.p[u] = v
		dsu.a[v] = x
		if dsu.r[u] == dsu.r[v] {
			dsu.r[v]++
		}
	}
}

type Query struct {
	from int
	to   int
}

func (dsu *DSU) DFS(g [][]int, q [][]Query, v int) {
	dsu.p[v] = v
	dsu.a[v] = v
	dsu.u[v] = true
	for _, t := range g[v] {
		if !dsu.u[t] {
			dsu.DFS(g, q, t)
			dsu.UnionSets(v, t, v)
		}
	}

	for _, qi := range q[v] {
		if dsu.u[qi.from] {
			qi.to = dsu.a[dsu.FindSet(qi.from)]
		}
	}
}

func LSA(g [][]int, q [][]Query) {
	a := NewDSU(len(g))
	a.DFS(g, q, 0)
}

func main() {
	var n int
	var l int
	fmt.Scan(&n)
	g := make([][]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&l)
		for j := 0; j < l; j++ {
			fmt.Scan(g[i][j])
		}
	}

	var m int
	fmt.Scan(&m)
	q := make([][]Query, n)
	for i := 0; i < m; i++ {
		var a int
		var b int
		fmt.Scan(&a, &b)
		q[a] = append(q[a], Query{from: b, to: -1})
		q[b] = append(q[b], Query{from: a, to: -1})
	}

	LSA(g, q)

	for i := 0; i < n; i++ {
		for j := 0; j < len(q[i]); j++ {
			fmt.Println("lsa(", i, ", ", q[i][j].from, ") = ", q[i][j].to)
		}
	}
}
