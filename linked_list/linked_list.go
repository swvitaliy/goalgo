package linked_list

// TODO add capacity (slice of Nodes, method Node.newNode)

type LinkedNode[T any] interface {
	Value() T
	Next() LinkedNode[T]
	Insert(value T)
	RemoveNext()
}

type LinkedList[T any] interface {
	Head() LinkedNode[T]
	Clear()
}

type Stack[T any] interface {
	PushFront(value T)
	PopFront()
	Empty() bool
	Top() (T, bool)
}

type Node[T any] struct {
	val  T
	next *Node[T]
}

func (n *Node[T]) Value() T {
	return n.val
}

func (n *Node[T]) Next() LinkedNode[T] {
	return n.next
}

func (n *Node[T]) Insert(value T) {
	n.next = &Node[T]{val: value, next: n.next}
}

func (n *Node[T]) RemoveNext() {
	if n.next != nil {
		n.next = n.next.next
	}
}

type List[T any] struct {
	head *Node[T]
}

func NewList[T any]() *List[T] {
	return &List[T]{
		head: &Node[T]{},
	}
}

func (l *List[T]) Head() LinkedNode[T] {
	return l.head.next
}

func (l *List[T]) Clear() {
	l.head = nil
}

/* Stack implementation */

func (l *List[T]) PushFront(value T) {
	l.head.Insert(value)
}

func (l *List[T]) PopFront() {
	l.head.RemoveNext()
}

func (l *List[T]) Top() (val T, ok bool) {
	ok = l.head != nil && l.head.next != nil
	if ok {
		val = l.head.next.val
	}
	return val, ok
}

func (l *List[T]) Empty() bool {
	return l.head == nil || l.head.next == nil
}
