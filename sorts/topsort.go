package sorts

import (
	"goalgo/slices"
)

func TopSort(graph [][]int) []int {
	n := len(graph)
	used := make([]bool, n)
	ans := make([]int, 0)
	var dfs func(int)
	dfs = func(v int) {
		used[v] = true
		for _, u := range graph[v] {
			if !used[u] {
				dfs(u)
			}
		}
		ans = append(ans, v)
	}
	for i := 0; i < n; i++ {
		if !used[i] {
			dfs(i)
		}
	}
	slices.Reverse(ans)
	return ans
}
