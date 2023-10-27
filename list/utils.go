package list

import (
	"cmp"
)

func IsSorted[T cmp.Ordered](a *Node[T], dir bool) bool {
	for aNext := a.Next; aNext != nil; aNext = aNext.Next {
		if dir == (a.Val > aNext.Val) {
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
	next := this.Next
	for next != nil {
		nextNext := next.Next
		next.Next = this
		this = next
		next = nextNext
	}
	head.Next = nil
	return this
}

func Reduce[T any](list *Node[T], acc func(res, v T) T, init T) T {
	res := init
	for list != nil {
		res = acc(res, list.Val)
		list = list.Next
	}
	return res
}
