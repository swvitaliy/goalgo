package main

import (
	"golang.org/x/exp/constraints"
	"lis/typ"
)

func Lis[T any, LessC typ.Comparator[T]](a []T, less LessC, b typ.Bounds[T]) []T {
	n := len(a)
	d := make([]T, n+1)
	for i := range d {
		d[i] = b.MaxValue()
	}
	d[0] = b.MinValue()

	prev := make([]int, n)
	pos := make([]int, n+1)
	pos[0] = -1

	l := 1
	for i, ai := range a {
		j := UpperBound(d, ai, less)
		if less.Compare(d[j-1], ai) && less.Compare(ai, d[j]) {
			d[j] = ai
			pos[j] = i
			prev[i] = pos[j-1]
			if j > l {
				l = j
			}
		}
	}

	var ans []T
	p := pos[l]
	for p != -1 {
		ans = append(ans, a[p])
		p = prev[p]
	}

	typ.Reverse(ans)
	return ans
}

func LisNumbers[T constraints.Signed](a []T, b typ.Bounds[T]) []T {
	var less = typ.LessOrdered[T]{}
	return Lis(a, less, b)
}

func LisNumbers_Example() {
	list := []int{1, 4, 6, 6, 32, 2, 7, 3, 34, 6, 6}
	LisNumbers(list, typ.IntBounds)
}
