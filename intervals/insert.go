package intervals

import "sort"

func Insert(a [][]int, x []int) [][]int {
	n := len(a)
	i := sort.Search(n, func(i int) bool { return a[i][1] >= x[0] })
	j := i
	for ; j < n && a[j][0] <= x[1]; j++ {
		x[0] = min(a[j][0], x[0])
		x[1] = max(a[j][1], x[1])
	}

	b := make([][]int, i+1+n-j)
	copy(b, a[:i])
	b[i] = x
	copy(b[i+1:], a[j:])
	return b
}
