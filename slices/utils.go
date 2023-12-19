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
	i := 0
	j := len(a) - 1
	for i < j {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}
}
func Reduce[S ~[]T, T any](list S, acc func(res, v T) T, init T) T {
	res := init
	for _, v := range list {
		res = acc(res, v)
	}
	return res
}

func Fill[S ~[]T, T any](a S, v T) {
	//for i := range a {
	//	a[i] = v
	//}
	// https://gist.github.com/taylorza/df2f89d5f9ab3ffd06865062a4cf015d
	a[0] = v
	for i := 1; i < len(a); i *= 2 {
		copy(a[i:], a[:i])
	}
}
