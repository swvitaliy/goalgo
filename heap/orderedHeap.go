package heap

import "cmp"

type OrderedHeap[T cmp.Ordered] []T

func NewOrderedHeap[T cmp.Ordered]() OrderedHeap[T] {
	return make(OrderedHeap[T], 0)
}

func (q *OrderedHeap[T]) sieveUp(i int) {
	a := *q
	for i > 0 {
		p := (i - 1) / 2
		if a[p] < a[i] {
			break
		}
		a[i], a[p] = a[p], a[i]
		i = p
	}
}

func (q *OrderedHeap[T]) sieveDown(i int) {
	a := *q
	n := len(a)

	for i < n {
		l := i*2 + 1
		r := i*2 + 2
		// min of 3 (i, l, r) should be parented
		toLeft := l < n && a[i] > a[l]
		toRight := r < n && a[i] > a[r]
		if toLeft && toRight {
			if a[l] < a[r] {
				a[i], a[l] = a[l], a[i]
				i = l
			} else {
				a[i], a[r] = a[r], a[i]
				i = r
			}
		} else if toLeft {
			a[i], a[l] = a[l], a[i]
			i = l
		} else if toRight {
			a[i], a[r] = a[r], a[i]
			i = r
		} else {
			break
		}
	}
}

func (q *OrderedHeap[T]) Insert(v T) {
	*q = append(*q, v)
	q.sieveUp(len(*q) - 1)
}

func (q *OrderedHeap[T]) Extract() T {
	a := *q
	n := len(a)
	if n == 0 {
		panic("cannot extract value from empty heap")
	}
	if n == 1 {
		r := a[0]
		*q = a[:0]
		return r
	}
	r := a[0]
	a[0] = a[n-1]
	*q = a[:n-1]
	q.sieveDown(0)
	return r
}

func (q *OrderedHeap[T]) Len() int {
	return len(*q)
}

func (q *OrderedHeap[T]) Peek() T {
	a := *q
	if len(a) == 0 {
		panic("cannot Peek value from empty heap")
	}
	return a[0]
}

func (q *OrderedHeap[T]) IsEmpty() bool {
	return len(*q) == 0
}
