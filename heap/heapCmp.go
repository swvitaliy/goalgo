package heap

import "cmp"

type comparator[T any] func(a, b T) bool

func Less[T cmp.Ordered](a, b T) bool {
	return a < b
}

func Greater[T cmp.Ordered](a, b T) bool {
	return a > b
}

type HeapCmp[T any, F comparator[T]] struct {
	data []T
	less F
}

func NewHeapCmp[T any, F comparator[T]](less F) HeapCmp[T, F] {
	return HeapCmp[T, F]{
		data: make([]T, 0),
		less: less,
	}
}

func (q *HeapCmp[T, F]) sieveUp(i int) {
	a := q.data
	for i > 0 {
		p := (i - 1) / 2
		if q.less(a[p], a[i]) {
			break
		}
		a[i], a[p] = a[p], a[i]
		i = p
	}
}

func (q *HeapCmp[T, F]) sieveDown(i int) {
	a := q.data
	n := len(a)

	for i < n {
		l := i*2 + 1
		r := i*2 + 2
		// min of 3 (i, l, r) should be parented
		toLeft := l < n && q.less(a[l], a[i])
		toRight := r < n && q.less(a[r], a[i])
		if toLeft && toRight {
			if q.less(a[l], a[r]) {
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

func (q *HeapCmp[T, F]) Len() int {
	return len(q.data)
}

func (q *HeapCmp[T, F]) IsEmpty() bool {
	return len(q.data) == 0
}

func (q *HeapCmp[T, F]) Peek() (T, bool) {
	a := q.data
	if len(a) == 0 {
		// panic("cannot Peek value from empty PriorQueue")
		var zero T
		return zero, false
	}
	return a[0], true
}

func (q *HeapCmp[T, F]) Enqueue(v T) {
	q.data = append(q.data, v)
	q.sieveUp(len(q.data) - 1)
}

func (q *HeapCmp[T, F]) Dequeue() (T, bool) {
	a := q.data
	n := len(a)
	if n == 0 {
		// panic("cannot Dequeue value from empty PriorQueue")
		var zero T
		return zero, false
	}
	if n == 1 {
		r := a[0]
		q.data = a[:0]
		return r, true
	}
	r := a[0]
	a[0] = a[n-1]
	q.data = a[:n-1]
	q.sieveDown(0)
	return r, true
}
