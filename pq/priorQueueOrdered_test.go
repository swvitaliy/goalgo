package pq

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestPriorQueueOrdered_Enqueue_Dequeue(t *testing.T) {
	q := NewPriorQueueOrdered[int]()
	q.Enqueue(0)
	q.Enqueue(22)
	q.Enqueue(7)
	assert.Equal(t, 0, q.Dequeue())
	assert.Equal(t, 7, q.Dequeue())
	assert.Equal(t, 22, q.Dequeue())
}

func TestPriorQueueOrdered1(t *testing.T) {
	a := []int{6, 9, 11, 15, 3, 8, 9, 2, 11, 7, 13, 12, 1, 9}
	q := NewPriorQueueOrdered[int]()
	for _, ai := range a {
		q.Enqueue(ai)
	}
	b := make([]int, 14)
	for i := 0; i < 14; i++ {
		b[i] = q.Dequeue()
	}
	sort.Ints(a)
	assert.Equal(t, a, b)
}
