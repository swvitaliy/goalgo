package sorts

import (
	"sync"
)

type sortFn = func(a, b []int, i, j int)

func merge(a, b []int, l, m, r int) {
	// b := make([]int, n-i+1, n-i+1)
	i := l
	k := l
	j := m + 1
	for i <= m && j <= r {
		if a[i] < a[j] {
			b[k] = a[i]
			i++
			k++
		} else {
			b[k] = a[j]
			j++
			k++
		}
	}
	if i <= m {
		copy(b[k:], a[i:m+1])
	}
	if j <= r {
		copy(b[k:], a[j:r+1])
	}
	copy(a[l:], b[l:r+1]) // copy(a[l:], b[:r-l+1]) when using a local b slice
}

// MergeSortParallel runs mergesort in parallel with k = log_2(number_of_goroutines)
// number_of_goroutines = 2^k
func MergeSortParallel(k int, a, b []int, i, j int, sort sortFn) {
	if k == 0 {
		sort(a, b, i, j)
	} else {
		m := (i + j) / 2
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			MergeSortParallel(k-1, a, b, i, m, sort)
		}()
		go func() {
			defer wg.Done()
			MergeSortParallel(k-1, a, b, m+1, j, sort)
		}()
		wg.Wait()
		merge(a, b, i, m, j)
	}
}

func insertionSort(a []int, i, j int) {
	for k := i + 1; k <= j; k++ {
		for l := k; l > i && a[l-1] > a[l]; l-- {
			a[l], a[l-1] = a[l-1], a[l]
		}
	}
}

// TimSort (switch to insertion sort for sub arrays of size 64 or less)
func TimSort(a, b []int, i, j int) {
	if j-i < 64 {
		insertionSort(a, i, j)
		return
	}

	m := (i + j) / 2
	TimSort(a, b, i, m)
	TimSort(a, b, m+1, j)
	merge(a, b, i, m, j)
}

func MergeSort(a, b []int, i, j int) {
	if i < j {
		m := (i + j) / 2
		MergeSort(a, b, i, m)
		MergeSort(a, b, m+1, j)
		merge(a, b, i, m, j)
	}
}
