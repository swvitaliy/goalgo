package algo

import (
	"goalgo/gen/gtyp"
	"golang.org/x/exp/constraints"
)

func Lis[T constraints.Signed](a []T, b gtyp.Bounds[T]) []T {
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
		j := UpperBound(d, func(v T) bool {
			return ai < v
		})
		if d[j-1] < ai && ai < d[j] {
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

	gtyp.Reverse(ans)
	return ans
}
