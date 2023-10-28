package graphs

import "fmt"

var g [][]int
var gr [][]int
var used []bool
var order []int
var comp []int

func dfs1(v int) {
	used[v] = true
	for _, u := range g[v] {
		if !used[u] {
			dfs1(u)
		}
	}
	order = append(order, v)
}

func dfs2(v int) {
	used[v] = true
	comp = append(comp, v)
	for _, u := range gr[v] {
		if !used[u] {
			dfs2(u)
		}
	}
}

func main() {
	var n int
	fmt.Scanf("%d", n)

	g = make([][]int, n)
	gr = make([][]int, n)
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Scanf("%d %d", a, b)
		g[a] = append(g[a], b)
		gr[b] = append(g[b], a)
	}

	used = make([]bool, n)
	for i := range used {
		used[i] = false
	}

	order = make([]int, n)

	for i := range g {
		if !used[i] {
			dfs1(i)
		}
	}

	for i := range used {
		used[i] = false
	}

	comp = make([]int, 0, n)

	for i := range g {
		v := order[n-i-1]
		if !used[v] {
			dfs2(v)
			comp = make([]int, 0, n)
		}
	}
}
