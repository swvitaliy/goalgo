package graphs

import "goalgo/slices"

// Scc for strongly connected components
func Scc(g, gr [][]int, n int, visit func(comp []int)) {
	var used []bool
	var order []int
	var comp []int

	var dfs func(v int)
	var dfsRev func(v int)

	dfs = func(v int) {
		used[v] = true
		for _, u := range g[v] {
			if !used[u] {
				dfs(u)
			}
		}
		order = append(order, v)
	}

	dfsRev = func(v int) {
		used[v] = true
		comp = append(comp, v)
		for _, u := range gr[v] {
			if !used[u] {
				dfsRev(u)
			}
		}
	}

	used = make([]bool, n)
	slices.Fill(used, false)

	order = make([]int, n)

	for i := range g {
		if !used[i] {
			dfs(i)
		}
	}

	slices.Fill(used, false)

	comp = make([]int, 0, n)

	for i := range g {
		v := order[n-i-1]
		if !used[v] {
			dfsRev(v)
			visit(comp)
			comp = make([]int, 0, n)
		}
	}
}
