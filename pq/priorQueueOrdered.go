package pq

import "cmp"

type PriorQueueOrdered[T cmp.Ordered] []T

func NewPriorQueueOrdered[T cmp.Ordered]() PriorQueueOrdered[T] {
	return make(PriorQueueOrdered[T], 0)
}

func (q *PriorQueueOrdered[T]) sieveUp(i int) {
	a := *q
	for i > 0 {
		p := i / 2
		if a[p] < a[i] {
			break
		}
		a[i], a[p] = a[p], a[i]
		i = p
	}
}

func (q *PriorQueueOrdered[T]) sieveDown(i int) {
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

func (q *PriorQueueOrdered[T]) Enqueue(v T) {
	*q = append(*q, v)
	q.sieveUp(len(*q) - 1)
}

func (q *PriorQueueOrdered[T]) Dequeue() T {
	a := *q
	n := len(a)
	if n == 0 {
		panic("cannot Dequeue value from empty PriorQueue")
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
