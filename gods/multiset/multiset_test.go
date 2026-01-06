package multiset

import (
	"testing"

	"github.com/emirpasic/gods/v2/trees/btree"
	"github.com/stretchr/testify/assert"
)

func mustKey[T any](x T, ok bool) T {
	if !ok {
		panic("Key is not in the multiset")
	}
	return x
}

func TestMsPut1(t *testing.T) {
	q := btree.New[int, int](3)
	MsPut(q, 1)
	MsPut(q, 2)
	MsPut(q, 3)
	MsPut(q, 4)
	MsPut(q, 5)
	assert.Equal(t, 5, q.Size())
}

func TestMsPut2(t *testing.T) {
	q := btree.New[int, int](3)
	MsPut(q, 1)
	MsPut(q, 2)
	MsPut(q, 3)
	MsPut(q, 1)
	MsPut(q, 2)
	assert.Equal(t, 3, q.Size())
	assert.Equal(t, 2, mustKey(q.Get(2)))
	assert.Equal(t, 2, mustKey(q.Get(2)))
}

func TestMsRemove1(t *testing.T) {
	q := btree.New[int, int](3)
	MsPut(q, 1)
	MsPut(q, 2)
	MsPut(q, 3)
	MsRemove(q, 1)
	assert.Equal(t, 2, q.Size())
	assert.Nil(t, q.GetNode(1))
}

func TestMsRemove2(t *testing.T) {
	q := btree.New[int, int](3)
	MsPut(q, 1)
	MsPut(q, 2)
	MsPut(q, 3)
	MsPut(q, 1)
	assert.Equal(t, 3, q.Size())
	assert.Equal(t, 2, mustKey(q.Get(1)))

	MsRemove(q, 1)
	assert.Equal(t, 3, q.Size())
	assert.Equal(t, 1, mustKey(q.Get(1)))

	MsRemove(q, 1)
	assert.Equal(t, 2, q.Size())
	assert.Nil(t, q.GetNode(1))
}
