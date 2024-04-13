package kmp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimple(t *testing.T) {
	p := "123"
	s := "123456"
	ans := SearchAll(p, s)
	assert.Equal(t, []int{2}, ans)
	assert.Equal(t, 1, SearchCount(p, s))
	assert.Equal(t, 2, SearchFirst(p, s))
}

func TestRepeatedPattern(t *testing.T) {
	p := "ababab"
	s := "abababababababababab"
	ans := SearchAll(p, s)
	expected := []int{0, 2, 4, 6, 8, 10, 12, 14}
	for i := range expected {
		expected[i] += len(p) - 1
	}
	assert.Equal(t, 8, SearchCount(p, s))
	assert.Equal(t, expected, ans)
}
