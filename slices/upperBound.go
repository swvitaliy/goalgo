package slices

import "cmp"

func UpperBound[S ~[]T, T cmp.Ordered](arr S, v T) int {
	l, r, m := 0, len(arr), 0

	for l < r {
		m = (l + r) / 2
		if cmp.Less(v, arr[m]) {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}
