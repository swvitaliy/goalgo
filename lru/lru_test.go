package lru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPutGet(t *testing.T) {
	lru := New[int, int](2)
	lru.Put(1, 1)
	_, ok := lru.Get(0)
	assert.Equal(t, false, ok)
	val, ok := lru.Get(1)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, val)
}

func TestOverflow(t *testing.T) {
	lru := New[int, int](2)
	lru.Put(1, 1)
	lru.Put(2, 2)
	lru.Put(3, 3)
	_, ok := lru.Get(1)
	assert.Equal(t, false, ok)
	_, ok = lru.Get(2)
	assert.Equal(t, true, ok)
	_, ok = lru.Get(3)
	assert.Equal(t, true, ok)
}
