package pq

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestPQueue_Enqueue_Dequeue(t *testing.T) {
	q := NewPriorQueuePair[uint, int]()
	q.Enqueue(0, 0)
	q.Enqueue(7, 7)
	q.Enqueue(2, 22)
	if v, _ := q.Dequeue(); v != 0 {
		t.Error()
	}
	if v, _ := q.Dequeue(); v != 22 {
		t.Error()
	}
	if v, _ := q.Dequeue(); v != 7 {
		t.Error()
	}
}
func TestPriorQueuePair1(t *testing.T) {
	a := []int{6, 9, 11, 15, 3, 8, 9, 2, 11, 7, 13, 12, 1, 9}
	q := NewPriorQueuePair[int, int]()
	for _, ai := range a {
		q.Enqueue(ai, ai)
	}
	b := make([]int, 14)
	for i := 0; i < 14; i++ {
		b[i], _ = q.Dequeue()
	}
	sort.Ints(a)
	assert.Equal(t, a, b)
}
