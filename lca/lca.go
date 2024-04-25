package lca

import (
	"cmp"
	"goalgo/trees/bst"
)

func BstLca[T cmp.Ordered](root *bst.Node[T], p, q T) *bst.Node[T] {
	t := root
	for t != nil {
		if p < t.Key && q < t.Key {
			t = t.Left
		} else if p > t.Key && q > t.Key {
			t = t.Right
		} else {
			return t
		}
	}
	return nil
}
