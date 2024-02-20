package rmq

import (
	"cmp"
)

type SegTree[T cmp.Ordered] []T

func NewSegTree[T cmp.Ordered](n int) SegTree[T] {
	return make([]T, 2*n)
}

// fn should be commutative
// fn should be monotonous
// fn should be associative
// fn should be idempotent
// For example:
// - sum
// - min
// - max
// - gcd
// - lcm
// - xor
// - or
// - and

// See
// https://www.geeksforgeeks.org/segment-tree-efficient-implementation/
// https://www.geeksforgeeks.org/iterative-segment-tree-range-maximum-query-with-node-update/

func BuildSegTree[T cmp.Ordered](a []T, fn func(a, b T) T) SegTree[T] {
	t := NewSegTree[T](len(a))
	n := len(a)
	for i := 0; i < n; i++ {
		t[n+i] = a[i]
	}

	for i := n - 1; i >= 0; i-- {
		t[i] = fn(t[i<<1], t[i<<1|1])
	}

	return t
}

func (t SegTree[T]) Update(pos int, newValue T, fn func(a, b T) T) {
	n := len(t) >> 1
	t[n+pos] = newValue

	for i := n + pos; i > 1; i >>= 1 {
		t[i>>1] = fn(t[i], t[i^1])
	}
}

func (t SegTree[T]) Query(l, r int, fn func(a, b T) T) T {
	r += 1

	n := len(t) >> 1
	var ans T
	l += n
	r += n
	for l < r {
		if l&1 != 0 {
			ans = fn(ans, t[l])
			l++
		}

		if r&1 != 0 {
			r--
			ans = fn(ans, t[r])
		}

		l >>= 1
		r >>= 1
	}

	return ans
}
