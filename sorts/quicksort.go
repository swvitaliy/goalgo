package sorts

func QuickSort(a []int, l, r int) {
	N := len(a)
	i := l
	j := r
	m := (l + r) / 2
	for i < j {
		for i < N && a[i] < a[m] {
			i++
		}
		for j >= 0 && a[j] > a[m] {
			j--
		}
		if i <= j {
			a[i], a[j] = a[j], a[i]
			i++
			j--
		}
	}
	if l < j {
		QuickSort(a, l, j)
	}
	if i < r {
		QuickSort(a, i, r)
	}
}
