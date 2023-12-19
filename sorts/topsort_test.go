package sorts

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTopSortEmpty(t *testing.T) {
	a := [][]int{}
	assert.Equal(t, []int{}, TopSort(a))
}

func TestTopSort0(t *testing.T) {
	a := [][]int{
		{1, 3},
		{0, 2, 4},
		{1, 5},
		{0, 6},
		{1, 5, 7},
		{2, 4, 8},
		{3, 7},
		{4, 6},
		{5, 9},
		{8},
	}
	assert.Equal(t, []int{0, 1, 2, 5, 8, 9, 4, 7, 6, 3}, TopSort(a))
}
