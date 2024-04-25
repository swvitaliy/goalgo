package lca

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

func GraphLca(g [][]int, q [][]Query) {
	a := NewDSU(len(g))
	a.DFS(g, q, 0)
}

/*
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

	LCA(g, q)

	for i := 0; i < n; i++ {
		for j := 0; j < len(q[i]); j++ {
			fmt.Println("lsa(", i, ", ", q[i][j].from, ") = ", q[i][j].to)
		}
	}
}
*/
