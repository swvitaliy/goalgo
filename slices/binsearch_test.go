package slices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinSearch(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, 2, BinSearch(a, 2), "BinSearch failed")
}

func TestBinSearch2(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, -1, BinSearch(a, 5), "BinSearch failed")
}

func TestBinSearch3(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, 0, BinSearch(a, 1), "BinSearch failed")
}

func TestBinSearch4(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, -1, BinSearch(a, 100500), "BinSearch failed")
}

func TestBinSearch5(t *testing.T) {
	a := []int{1, 3, 2, 7, 0}
	assert.Equal(t, -1, BinSearch(a, -10), "BinSearch failed")
}
