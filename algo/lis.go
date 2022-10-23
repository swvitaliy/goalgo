package algo

import (
	"golang.org/x/exp/constraints"
	"math"
	"sort"
)

func Lis[T constraints.Signed](a []T) []T {
	n := len(a)
	d := make([]T, n+1)
	for i := range d {
		d[i] = T(math.Inf(1))
	}
	d[0] = T(math.Inf(-1))

	prev := make([]T, n)
	pos := make([]T, n+1)
	pos[0] = -1

	l := 1
	for i, _ := range a {
		j := sort.Search(n, func(k int) bool {
			return d[k] == a[i]
		})
		if d[j-1] < a[i] && a[i] < d[j] {
			d[j] = a[i]
			pos[j] = T(i)
			prev[i] = pos[j-1]
			if j > l {
				l = j
			}
		}
	}

	var ans []T
	p := pos[l]
	for p != -1 {
		ans = append(ans, a[p])
		p = prev[p]
	}

	return ans
}
