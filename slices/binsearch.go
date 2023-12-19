package slices

import "cmp"

func BinSearch[S ~[]T, T cmp.Ordered](a S, v T) int {
	i, j := 0, len(a)-1
	for i <= j {
		m := (i + j) / 2
		if a[m] == v {
			return m
		}
		if v < a[m] {
			j = m - 1
		} else {
			i = m + 1
		}
	}
	return -1
}
