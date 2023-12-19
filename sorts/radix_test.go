package sorts

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRadixSortEmpty(t *testing.T) {
	a := []int{}
	assert.Equal(t, []int{}, RadixSort(a))
}
func TestRadixSortOne(t *testing.T) {
	a := []int{777}
	assert.Equal(t, []int{777}, RadixSort(a))
}

func TestRadixSort0(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, RadixSort(a))
}
func TestRadixSort1(t *testing.T) {
	a := []int{5, 4, 3, 2, 1}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, RadixSort(a))
}

func TestRadixSort2(t *testing.T) {
	a := []int{1, 4, 6, 6, 32, 2, 7, 3, 34, 6, 6}
	assert.Equal(t, []int{1, 2, 3, 4, 6, 6, 6, 6, 7, 32, 34}, RadixSort(a))
}

func TestRadixSort3(t *testing.T) {
	a := []int{654, 100, 572, 938, 807, 587, 951, 395, 794, 124, 758, 271, 828, 522, 876, 750, 657, 165, 368, 423}
	assert.Equal(t, []int{100, 124, 165, 271, 368, 395, 423, 522, 572, 587, 654, 657, 750, 758, 794, 807, 828, 876, 938, 951}, RadixSort(a))
}
