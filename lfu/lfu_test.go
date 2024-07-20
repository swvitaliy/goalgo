package lfu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLFUSimple(t *testing.T) {
	cache := New[int, int](2, -1)
	cache.Put(1, 1)
	cache.Get(1)
	cache.Put(2, 2)
	assert.Equal(t, 1, cache.Get(1))
	assert.Equal(t, 2, cache.Get(2))
	cache.Put(3, 3)
	assert.Equal(t, 3, cache.Get(3))
	assert.Equal(t, -1, cache.Get(2))
	assert.Equal(t, 1, cache.Get(1))
}
