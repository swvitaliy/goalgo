package heap

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestPairHeap_Insert_Extract(t *testing.T) {
	q := NewPairHeap[uint, int]()
	q.Insert(0, 3456)
	q.Insert(7, 8456)
	q.Insert(2, 782)
	assert.Equal(t, uint(0), Priority(q.Peek()))
	assert.Equal(t, 3456, Value(q.Extract()))
	assert.Equal(t, uint(2), Priority(q.Peek()))
	assert.Equal(t, 782, Value(q.Extract()))
	assert.Equal(t, uint(7), Priority(q.Peek()))
	assert.Equal(t, 8456, Value(q.Extract()))
}
func TestPairHeap1(t *testing.T) {
	a := []int{6, 9, 11, 15, 3, 8, 9, 2, 11, 7, 13, 12, 1, 9}
	q := NewPairHeap[int, int]()
	for _, ai := range a {
		q.Insert(ai, ai)
	}
	b := make([]int, 14)
	for i := 0; i < 14; i++ {
		b[i], _ = q.Extract()
	}
	sort.Ints(a)
	assert.Equal(t, a, b)
}
