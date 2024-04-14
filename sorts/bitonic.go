package sorts

import (
	"cmp"
	"goalgo/limits"
	"goalgo/slices"
	"math/bits"
)

const Ascending = true
const Descending = false

func BitonicSort[T cmp.Ordered](a []T, dir bool) []T {
	m := len(a)
	a = bitonicFillUntil2k(a, dir)
	n := len(a)
	if n <= 1 {
		return a
	}
	k := n / 2
	BitonicSort(a[:k], Ascending)
	BitonicSort(a[k:], Descending)
	bitonicMerge(a, dir)
	return a[:m]
}

func bitonicFillUntil2k[T cmp.Ordered](a []T, dir bool) []T {
	n := len(a)
	k := 63 - bits.LeadingZeros64(uint64(n)) // lg(n)
	m := 1 << k
	if n == m {
		return a
	}
	m = 1 << (k + 1)
	a = append(a, make([]T, m-n)...)
	if dir == Ascending {
		slices.Fill(a[n:], limits.MaxValue[T]())
	} else {
		slices.Fill(a[n:], limits.MinValue[T]())
	}
	return a
}

func bitonicMerge[T cmp.Ordered](a []T, dir bool) {
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
	bitonicMerge(a[:k], dir)
	bitonicMerge(a[k:], dir)
}
