package lca

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
		ans.p[i] = i
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
