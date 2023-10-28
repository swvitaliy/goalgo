package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverse(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	assert.True(t, IsSorted(a, true), "List is not sorted (a-z)")
	Reverse(a)
	assert.True(t, IsSorted(a, false), "Reversed is not sorted (z-a)")
}

func TestReduce1(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	assert.Equal(t, Reduce(a, func(res, v int) int { return res + v }, 0), 15, "Reduce failed")
}
