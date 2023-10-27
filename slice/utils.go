package slice

import "cmp"

func IsSorted[T cmp.Ordered](a []T, dir bool) bool {
	for i := 1; i < len(a); i++ {
		if dir == (a[i-1] > a[i]) {
			return false
		}
	}
	return true
}

func Reverse[T any](a []T) {
	var n = len(a)
	var m = n >> 1
	for i := 0; i < m; i++ {
		j := n - i - 1
		a[i], a[j] = a[j], a[i]
	}
}
func Reduce[T any](list []T, acc func(res, v T) T, init T) T {
	res := init
	for _, v := range list {
		res = acc(res, v)
	}
	return res
}
