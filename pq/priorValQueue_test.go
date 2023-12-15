package pq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPriorValQueue_Enqueue_Dequeue(t *testing.T) {
	q := NewPriorValQueue[uint]()
	q.Enqueue(0)
	q.Enqueue(22)
	q.Enqueue(7)
	assert.Equal(t, 0, q.Dequeue())
	assert.Equal(t, 7, q.Dequeue())
	assert.Equal(t, 22, q.Dequeue())
}
