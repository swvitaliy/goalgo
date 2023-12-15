package graphs

import "goalgo/slices"

type Color int

const (
	White Color = iota
	Gray
	Black
)

// CycleDFS for cycle detection
// Returns index of some node of the cycle
// Returns -1 if there is no cycle
func CycleDFS(g [][]int) int {
	n := len(g)
	c := make([]Color, n)
	p := make([]int, n)
	slices.Fill(c, White)

	for i := 0; i < n; i++ {
		t := cycleDFS(g, c, p, i)
		if t != -1 {
			return t
		}
	}

	return -1
}

func CyclesDFS(g [][]int) []int {
	n := len(g)
	c := make([]Color, n)
	p := make([]int, n)
	slices.Fill(c, White)

	var ans []int
	for i := 0; i < n; i++ {
		t := cycleDFS(g, c, p, i)
		if t != -1 {
			ans = append(ans, t)
		}
	}

	return ans
}

func cycleDFS(g [][]int, c []Color, p []int, i int) int {
	c[i] = Gray
	for _, u := range g[i] {
		if c[u] == Gray {
			return u
		}
		if c[u] == White {
			p[u] = i
			if cycleDFS(g, c, p, u) != -1 {
				return u
			}
		}
	}

	return -1
}
