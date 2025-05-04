package linked_list

import (
	"cmp"
)

func IsSorted[T cmp.Ordered](a *Node[T], dir bool) bool {
	for aNext := a.next; aNext != nil; aNext = aNext.next {
		if dir == (a.val > aNext.val) {
			return false
		}
		a = aNext
	}
	return true
}

func Reverse[T any](head *Node[T]) *Node[T] {
	if head == nil {
		return nil
	}
	this := head
	next := this.next
	for next != nil {
		nextNext := next.next
		next.next = this
		this = next
		next = nextNext
	}
	head.next = nil
	return this
}

func Reduce[T any](list *Node[T], acc func(res, v T) T, init T) T {
	res := init
	for list != nil {
		res = acc(res, list.val)
		list = list.next
	}
	return res
}
