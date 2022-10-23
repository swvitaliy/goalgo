package main

import "lis/typ"

func UpperBound[T any, LessC typ.Comparator[T]](arr []T, x T, less LessC) int {
	l, r, m := 0, len(arr), 0

	for l < r {
		m = (l + r) / 2
		if less.Compare(x, arr[m]) {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}
