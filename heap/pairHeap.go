package heap

import "cmp"

// min-heap
type pair[P cmp.Ordered, T any] struct {
	prior P
	val   T
}

type PairHeap[P cmp.Ordered, T any] []pair[P, T]

func NewPairHeap[P cmp.Ordered, T any]() PairHeap[P, T] {
	return make(PairHeap[P, T], 0)
}

func Priority[P any, T any](p P, v T) P {
	return p
}

func Value[P any, T any](p P, v T) T {
	return v
}

func (q *PairHeap[P, T]) sieveUp(i int) {
	a := *q
	for i > 0 {
		p := (i - 1) / 2
		if a[p].prior < a[i].prior {
			break
		}
		a[i], a[p] = a[p], a[i]
		i = p
	}
}

func (q *PairHeap[P, T]) sieveDown(i int) {
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

func (q *PairHeap[P, T]) Insert(p P, v T) {
	*q = append(*q, pair[P, T]{p, v})
	q.sieveUp(len(*q) - 1)
}

func (q *PairHeap[P, T]) Extract() (P, T) {
	a := *q
	n := len(a)
	if n == 0 {
		panic("cannot extract value from empty PairHeap")
	}
	if n == 1 {
		r := a[0]
		*q = a[:0]
		return r.prior, r.val
	}
	r := a[0]
	a[0] = a[n-1]
	*q = a[:n-1]
	q.sieveDown(0)
	return r.prior, r.val
}
func (q *PairHeap[P, T]) Len() int {
	return len(*q)
}

func (q *PairHeap[P, T]) Peek() (P, T) {
	a := *q
	if len(a) == 0 {
		panic("cannot Peek value from empty PairHeap")
	}
	return a[0].prior, a[0].val
}

func (q *PairHeap[P, T]) IsEmpty() bool {
	return len(*q) == 0
}
