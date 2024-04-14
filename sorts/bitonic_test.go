package sorts

import (
	"github.com/stretchr/testify/assert"
	"goalgo/slices"
	"testing"
)

func TestBitonicSort(t *testing.T) {
	a := []int{1, 4, 6, 6, 32, 2, 7, 3, 34, 6, 6}
	aSorted := []int{1, 2, 3, 4, 6, 6, 6, 6, 7, 32, 34}
	a = BitonicSort(a, Ascending)
	assert.True(t, slices.IsSorted(a, slices.Ascending))
	assert.Equal(t, aSorted, a)

}
func TestBitonicSortDescending(t *testing.T) {
	a := []int{1, 4, 6, 6, 32, 2, 7, 3, 34, 6, 6}
	aSorted := []int{1, 2, 3, 4, 6, 6, 6, 6, 7, 32, 34}
	slices.Reverse(aSorted)
	a = BitonicSort(a, Descending)
	assert.True(t, slices.IsSorted(a, slices.Descending))
	assert.Equal(t, aSorted, a)
}

func TestBitonicInplaceWhen2K(t *testing.T) {
	a := []int{1, 4, 6, 6, 32, 2, 7, 3}
	b := BitonicSort(a, Ascending)
	assert.Equal(t, &b, &a)
}
