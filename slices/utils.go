package slices

import "cmp"

const Ascending = true
const Descending = false

func IsSorted[S ~[]T, T cmp.Ordered](a S, dir bool) bool {
	for i := 1; i < len(a); i++ {
		if dir == (a[i-1] > a[i]) {
			return false
		}
	}
	return true
}

func Reverse[S ~[]T, T any](a S) {
	var n = len(a)
	var m = n >> 1
	for i := 0; i < m; i++ {
		j := n - i - 1
		a[i], a[j] = a[j], a[i]
	}
}
func Reduce[S ~[]T, T any](list S, acc func(res, v T) T, init T) T {
	res := init
	for _, v := range list {
		res = acc(res, v)
	}
	return res
}