package goalgo

// min-heap
type pqItem[T any] struct {
	prior int
	val   T
}

type PQueue[T any] []pqItem[T]

func NewPQueue[T any]() PQueue[T] {
	return make(PQueue[T], 0)
}

func (q *PQueue[T]) sieveUp(i int) {
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

func (q *PQueue[T]) sieveDown(i int) {
	a := *q
	n := len(a)

	swp := func(i, l int) {
		a[i], a[l] = a[l], a[i]
		i = l
	}

	for i < n {
		l := i*2 + 1
		r := i*2 + 2
		// min of 3 (i, l, r) should be parented
		toLeft := l < n && a[i].prior > a[l].prior
		toRight := r < n && a[i].prior > a[r].prior
		if toLeft && toRight {
			if a[l].prior < a[r].prior {
				swp(i, l)
			} else {
				swp(i, r)
			}
		} else if toLeft {
			swp(i, l)
		} else if toRight {
			swp(i, r)
		} else {
			break
		}
	}
}

func (q *PQueue[T]) Enqueue(p int, v T) {
	*q = append(*q, pqItem[T]{p, v})
	q.sieveUp(len(*q) - 1)
}

func (q *PQueue[T]) Dequeue() (T, int) {
	a := *q
	n := len(a)
	if n == 0 {
		panic("cannot Dequeue value from empty PQueue")
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
