package sorts

func QuickSort(a []int, l, r int) {
	i := l
	j := r
	pivot := a[(l+r)/2]
	for i < j {
		for a[i] < pivot {
			i++
		}
		for a[j] > pivot {
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
