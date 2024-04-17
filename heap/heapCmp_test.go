package heap

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestPriorQueueCmp_Enqueue_Dequeue(t *testing.T) {
	q := NewHeapCmp[uint, comparator[uint]](Less[uint])
	q.Enqueue(0)
	q.Enqueue(22)
	q.Enqueue(7)
	assert.Equal(t, uint(0), q.Dequeue())
	assert.Equal(t, uint(7), q.Dequeue())
	assert.Equal(t, uint(22), q.Dequeue())
}

func TestPriorQueueCmp1(t *testing.T) {
	a := []int{6, 9, 11, 15, 3, 8, 9, 2, 11, 7, 13, 12, 1, 9}
	q := NewHeapCmp[int, comparator[int]](Less[int])
	for _, ai := range a {
		q.Enqueue(ai)
	}
	assert.Equal(t, 14, q.Len())
	b := make([]int, 14)
	for i := 0; i < 14; i++ {
		b[i] = q.Dequeue()
	}
	sort.Ints(a)
	assert.Equal(t, a, b)
}
