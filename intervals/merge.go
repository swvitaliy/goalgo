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
		e := r[k-1] // e is the last element in the result
		if a[i][0] <= e[1] {
			e[1] = max(e[1], a[i][1])
		} else {
			r = append(r, a[i])
			k++
		}
	}

	return r
}
