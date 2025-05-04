package linked_list

import "cmp"

// TODO split Find method into lower/upper bounds

type SortedList[T any] interface {
	LinkedList[T]
	FindFn(target T, cmp func(a, b T) int) LinkedNode[T]
	Contains(value T, cmp func(a, b T) int) bool
}

func (l *List[T]) FindFn(target T, cmp func(a, b T) int) LinkedNode[T] {
	n := l.head
	if n == nil || n.next == nil {
		return nil
	}

	n = n.next
	for ; n != nil; n = n.next {
		c := cmp(n.val, target)
		if c < 0 || c == 0 {
			return n
		}
	}

	return nil
}

func (l *List[T]) Contains(value T, cmp func(a, b T) int) bool {
	n := l.FindFn(value, cmp)
	if n == nil {
		return false
	}
	return cmp(n.Value(), value) == 0
}

type OrderedList[T cmp.Ordered] interface {
	LinkedList[T]
	Find(target T) LinkedNode[T]
	Contains(value T) bool
}

type orderedList[T cmp.Ordered] struct {
	*List[T]
}

func (l *orderedList[T]) Head() LinkedNode[T] {
	return l.List.Head()
}

func (l *orderedList[T]) Clear() {
	l.List.Clear()
}

func NewOrderedList[T cmp.Ordered]() OrderedList[T] {
	return &orderedList[T]{
		List: NewList[T](),
	}
}

func (l *orderedList[T]) Find(target T) LinkedNode[T] {
	n := l.head
	if n == nil || n.next == nil {
		return nil
	}

	n = n.next
	for ; n != nil; n = n.next {
		if n.val < target || n.val == target {
			return n
		}
	}

	return nil
}

func (l *orderedList[T]) Contains(value T) bool {
	n := l.Find(value)
	if n == nil {
		return false
	}

	return n.Value() == value
}
