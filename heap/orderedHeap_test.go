package heap

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestOrderedHeap_Insert_Extract(t *testing.T) {
	q := NewOrderedHeap[int]()
	q.Insert(0)
	q.Insert(22)
	q.Insert(7)
	assert.Equal(t, 0, q.Extract())
	assert.Equal(t, 7, q.Extract())
	assert.Equal(t, 22, q.Extract())
}

func TestPriorQueueOrdered1(t *testing.T) {
	a := []int{6, 9, 11, 15, 3, 8, 9, 2, 11, 7, 13, 12, 1, 9}
	q := NewOrderedHeap[int]()
	for _, ai := range a {
		q.Insert(ai)
	}
	b := make([]int, 14)
	for i := 0; i < 14; i++ {
		b[i] = q.Extract()
	}
	sort.Ints(a)
	assert.Equal(t, a, b)
}
