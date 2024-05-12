package heap

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestHeap_Insert_Extract(t *testing.T) {
	q := NewHeap[uint, comparator[uint]](Less[uint])
	q.Insert(0)
	q.Insert(22)
	q.Insert(7)
	assert.Equal(t, uint(0), q.Extract())
	assert.Equal(t, uint(7), q.Extract())
	assert.Equal(t, uint(22), q.Extract())
}

func TestHeap_Insert_Extract_Max(t *testing.T) {
	q := NewHeap[uint, comparator[uint]](Greater[uint])
	q.Insert(0)
	q.Insert(22)
	q.Insert(7)
	assert.Equal(t, uint(22), q.Extract())
	assert.Equal(t, uint(7), q.Extract())
	assert.Equal(t, uint(0), q.Extract())
}

func TestHeap1(t *testing.T) {
	a := []int{6, 9, 11, 15, 3, 8, 9, 2, 11, 7, 13, 12, 1, 9}
	q := NewHeap[int, comparator[int]](Less[int])
	for _, ai := range a {
		q.Insert(ai)
	}
	assert.Equal(t, 14, q.Len())
	b := make([]int, 14)
	for i := 0; i < 14; i++ {
		b[i] = q.Extract()
	}
	sort.Ints(a)
	assert.Equal(t, a, b)
}
