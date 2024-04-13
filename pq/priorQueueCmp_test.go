package pq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPriorQueueCmp_Enqueue_Dequeue(t *testing.T) {
	q := NewPriorQueueCmp[uint, comparator[uint]](Less[uint])
	q.Enqueue(0)
	q.Enqueue(22)
	q.Enqueue(7)
	assert.Equal(t, uint(0), q.Dequeue())
	assert.Equal(t, uint(7), q.Dequeue())
	assert.Equal(t, uint(22), q.Dequeue())
}
