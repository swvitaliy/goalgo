package bst

import "cmp"

type Node[T cmp.Ordered] struct {
	Key   T
	Left  *Node[T]
	Right *Node[T]
}

func MakeFromLinear[T cmp.Ordered](a []T, zero T) *Node[T] {
	if len(a) == 0 {
		return nil
	}

	if len(a) >= 1 && a[0] == zero {
		return nil
	}

	var left, right *Node[T]
	if len(a) > 1 {
		left = MakeFromLinear[T](a[1:], zero)
	}
	if len(a) > 2 {
		right = MakeFromLinear[T](a[2:], zero)
	}

	return &Node[T]{Key: a[0], Left: left, Right: right}
}

func Search[T cmp.Ordered](t *Node[T], key T) *Node[T] {
	if t == nil {
		return nil
	}

	for t != nil && t.Key != key {
		if t.Key < key {
			t = t.Right
		} else {
			t = t.Left
		}
	}

	return t
}
