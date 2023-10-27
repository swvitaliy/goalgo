package sorts

import "cmp"

func BitonicSort[T cmp.Ordered](a []T, dir bool) {
	n := len(a)
	if n == 1 {
		return
	}
	k := n / 2
	BitonicSort(a[:k], dir)
	BitonicSort(a[k:], !dir)
	merge(a, dir)
}

func merge[T cmp.Ordered](a []T, dir bool) {
	n := len(a)
	if n <= 1 {
		return
	}
	k := n / 2
	for i := 0; i < k; i++ {
		if dir == (a[i] > a[i+k]) {
			a[i], a[i+k] = a[i+k], a[i]
		}
	}
	merge(a[:k], dir)
	merge(a[k:], dir)
}
