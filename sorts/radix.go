package sorts

import (
	"slices"
)

type RadixType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func RadixSort[T RadixType](a []T) []T {
	if len(a) < 2 {
		return a
	}
	b := make([]T, len(a))
	c := make([]int, 10)
	mx := slices.Max(a)
	// how many digits in the largest number (mx)?
	for exp := T(1); mx/exp > 0; exp *= 10 {
		a, b = CountingSort(a, b, c, exp), a
	}

	return a
}

func CountingSort[T RadixType](a, b []T, c []int, exp T) []T {
	for i := range c {
		c[i] = 0
	}
	for i := range a {
		c[(a[i]/exp)%10]++
	}
	for i := 1; i < 10; i++ {
		c[i] += c[i-1]
	}
	for i := len(a) - 1; i >= 0; i-- {
		t := a[i] / exp
		b[c[t%10]-1] = a[i]
		c[t%10]--
	}
	return b
}
