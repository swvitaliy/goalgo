package slices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinSearch(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, 2, BinSearch(a, 2), "BinSearch failed")
}

func TestBinSearchN(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	cmp := func(v int) func(int) int {
		return func(i int) int {
			if a[i] == v {
				return 0
			} else if a[i] < v {
				return -1
			}
			return 1
		}
	}
	assert.Equal(t, 2, BinSearchN(len(a), cmp(2)), "BinSearch failed")
}

func TestBinSearch2(t *testing.T) {
	a := []int{0, 1, 2, 3, 7}
	assert.Equal(t, -1, BinSearch(a, 5), "BinSearch failed")
}

func TestBinSearchN2(t *testing.T) {
	a := []int{0, 1, 2, 3, 7}
	cmp := func(v int) func(int) int {
		return func(i int) int {
			if a[i] == v {
				return 0
			} else if a[i] < v {
				return -1
			}
			return 1
		}
	}
	assert.Equal(t, -1, BinSearchN(len(a), cmp(5)), "BinSearch failed")
}

func TestBinSearch3(t *testing.T) {
	a := []int{0, 1, 2, 3, 7}
	assert.Equal(t, 1, BinSearch(a, 1), "BinSearch failed")
}

func TestBinSearch4(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, -1, BinSearch(a, 100500), "BinSearch failed")
}

func TestBinSearch5(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, -1, BinSearch(a, -10), "BinSearch failed")
}

func TestUpperLower1(t *testing.T) {
	a := []int{1, 7, 23, 56, 67}
	assert.Equal(t, 2, LowerBound(a, 23), "LowerBound failed")
	assert.Equal(t, 3, UpperBound(a, 23), "UpperBound failed")
}

func TestUpperLower2(t *testing.T) {
	a := []int{1, 7, 23, 56, 67}
	assert.Equal(t, 2, LowerBound(a, 17), "LowerBound failed")
	assert.Equal(t, 2, UpperBound(a, 17), "UpperBound failed")
}

func TestUpperLowerN1(t *testing.T) {
	a := []int{1, 7, 23, 56, 67}
	cmp := func(v int) func(int) int {
		return func(i int) int {
			if a[i] == v {
				return 0
			}
			if a[i] < v {
				return -1
			}
			return 1
		}
	}
	assert.Equal(t, 2, LowerBoundN(len(a), cmp(23)), "LowerBound failed")
	assert.Equal(t, 3, UpperBoundN(len(a), cmp(23)), "UpperBound failed")
}

func TestUpperLowerN2(t *testing.T) {
	a := []int{1, 7, 23, 56, 67}
	cmp := func(v int) func(int) int {
		return func(i int) int {
			if a[i] == v {
				return 0
			}
			if a[i] < v {
				return -1
			}
			return 1
		}
	}
	assert.Equal(t, 2, UpperBoundN(len(a), cmp(17)), "UpperBound failed")
	assert.Equal(t, 2, LowerBoundN(len(a), cmp(17)), "LowerBound failed")
}
