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

func BinSearchN(n int, cmp func(int) int) int {
	i, j := 0, n
	for i <= j {
		m := (i + j) / 2
		compareResult := cmp(m)
		if compareResult == 0 {
			return m
		}
		if compareResult > 0 {
			i = m + 1
		} else {
			j = m - 1
		}
	}
	return -1
}

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

func UpperBoundN(n int, cmp func(int) int) int {
	l, r, m := 0, n, 0
	for l < r {
		m = (l + r) / 2
		if cmp(m) == 1 { // less than
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}

func LowerBound[S ~[]T, T cmp.Ordered](arr S, v T) int {
	l, r, m := 0, len(arr), 0

	for l < r {
		m = (l + r) / 2
		if cmp.Less(arr[m], v) {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func LowerBoundN(n int, cmp func(int) int) int {
	l, r, m := 0, n, 0
	for l < r {
		m = (l + r) / 2
		if cmp(m) == -1 {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}
