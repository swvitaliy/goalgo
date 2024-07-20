package multiset

import (
	"github.com/emirpasic/gods/v2/trees/btree"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type Must2Func[T any] func(T, bool) T

func MakeMust2Func[T any](pattern string) Must2Func[T] {
	return func(x T, ok bool) T {
		if !ok {
			log.Fatalf(pattern)
		}
		return x
	}
}

var Must = MakeMust2Func[int]("Key is not in the multiset")

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
	assert.Equal(t, 2, Must(q.Get(2)))
	assert.Equal(t, 2, Must(q.Get(2)))
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
	assert.Equal(t, 2, Must(q.Get(1)))

	MsRemove(q, 1)
	assert.Equal(t, 3, q.Size())
	assert.Equal(t, 1, Must(q.Get(1)))

	MsRemove(q, 1)
	assert.Equal(t, 2, q.Size())
	assert.Nil(t, q.GetNode(1))
}
