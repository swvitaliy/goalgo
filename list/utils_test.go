package list

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverse(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	list := FromSlice(a)
	assert.True(t, IsSorted(list, true), "List is not sorted (a-z)")
	reversed := Reverse(list)
	assert.True(t, IsSorted(reversed, false), "Reversed is not sorted (z-a)")
	assert.Equal(t, []int{5, 4, 3, 2, 1}, ToSlice(reversed), "Reverse failed")
}

func TestReduce1(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	assert.Equal(t, Reduce(FromSlice(a), func(res, v int) int { return res + v }, 0), 15, "Reduce failed")
}
