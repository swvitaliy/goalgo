package pq

import "cmp"

// min-heap
type pqItem[T any, P cmp.Ordered] struct {
	prior P
	val   T
}

type PriorQueue[T any, P cmp.Ordered] []pqItem[T, P]

func NewPriorQueue[T any, P cmp.Ordered]() PriorQueue[T, P] {
	return make(PriorQueue[T, P], 0)
}

func NewPQ[T any]() PriorQueue[T, int] {
	return make(PriorQueue[T, int], 0)
}

func (q *PriorQueue[T, P]) sieveUp(i int) {
	a := *q
	for i > 0 {
		p := i / 2
		if a[p].prior < a[i].prior {
			break
		}
		a[i], a[p] = a[p], a[i]
		i = p
	}
}

func (q *PriorQueue[T, P]) sieveDown(i int) {
	a := *q
	n := len(a)

	for i < n {
		l := i*2 + 1
		r := i*2 + 2
		// min of 3 (i, l, r) should be parented
		toLeft := l < n && a[i].prior > a[l].prior
		toRight := r < n && a[i].prior > a[r].prior
		if toLeft && toRight {
			if a[l].prior < a[r].prior {
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

func (q *PriorQueue[T, P]) Enqueue(p P, v T) {
	*q = append(*q, pqItem[T, P]{p, v})
	q.sieveUp(len(*q) - 1)
}

func (q *PriorQueue[T, P]) Dequeue() (T, P) {
	a := *q
	n := len(a)
	if n == 0 {
		panic("cannot Dequeue value from empty PriorQueue")
	}
	if n == 1 {
		r := a[0]
		*q = a[:0]
		return r.val, r.prior
	}
	r := a[0]
	a[0] = a[n-1]
	*q = a[:n-1]
	q.sieveDown(0)
	return r.val, r.prior
}
