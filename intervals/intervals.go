package intervals

import "sort"

func Merge(a [][]int) [][]int {
	n := len(a)
	if n == 0 {
		return [][]int{}
	}
	sort.Slice(a, func(i, j int) bool { return a[i][0] < a[j][0] })
	r := make([][]int, 0)
	r = append(r, a[0])
	k := 1
	for i := 1; i < n; i++ {
		e := r[k-1] // e is last element in result
		if a[i][0] <= e[1] {
			e[1] = max(e[1], a[i][1])
		} else {
			r = append(r, a[i])
			k++
		}
	}

	return r
}

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
