package algo

type Comparator[T any] func(a T) bool

func UpperBound[T any](arr []T, less Comparator[T]) int {
	l, r, m := 0, len(arr), 0

	for l < r {
		m = (l + r) / 2
		if less(arr[m]) {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}
