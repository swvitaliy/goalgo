package algo

func Reverse[T any](a []T) {
	var n = len(a)
	var m = n >> 1
	for i := 0; i < m; i++ {
		j := n - i - 1
		a[i], a[j] = a[j], a[i]
	}
}
