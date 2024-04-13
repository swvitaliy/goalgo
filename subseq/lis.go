package subseq

import (
	"cmp"
	"goalgo/limits"
	"goalgo/slices"
)

func Lis[S ~[]T, T cmp.Ordered](a S) S {
	n := len(a)
	d := make([]T, n+1)
	for i := range d {
		d[i] = limits.MaxValue[T]()
	}
	d[0] = limits.MinValue[T]()

	prev := make([]int, n)
	pos := make([]int, n+1)
	pos[0] = -1

	l := 1
	for i, ai := range a {
		j := slices.UpperBound(d, ai)
		if cmp.Less(d[j-1], ai) && cmp.Less(ai, d[j]) {
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

	slices.Reverse(ans)
	return ans
}
