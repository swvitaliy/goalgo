package rbtree

import "cmp"

type Color bool

const (
	Red   Color = true
	Black Color = false
)

type Node[T cmp.Ordered] struct {
	Key    T
	Color  Color
	Parent *Node[T]
	Left   *Node[T]
	Right  *Node[T]
}

func NewRBNode[T cmp.Ordered](key T, color Color, left *Node[T], right *Node[T]) *Node[T] {
	return &Node[T]{Key: key, Color: color, Left: left, Right: right}
}

func (t *Node[T]) Insert(key T) *Node[T] {
	x := NewRBNode(key, Red, nil, nil)
	if t == nil {
		x.Parent = nil
		return fixInsertion(x)
	}

	p := t
	var q *Node[T] = nil
	for p != nil {
		q = p
		if p.Key < t.Key {
			p = p.Right
		} else {
			p = p.Left
		}
	}

	x.Parent = p
	if q.Key < t.Key {
		q.Right = t
	} else {
		q.Left = t
	}

	return fixInsertion(t)
}

func fixInsertion[T cmp.Ordered](t *Node[T]) *Node[T] {
	if t.Parent == nil {
		t.Color = Black
	}

}
