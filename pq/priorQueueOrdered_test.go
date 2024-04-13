package pq

import (
	"github.com/stretchr/testify/assert"
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
