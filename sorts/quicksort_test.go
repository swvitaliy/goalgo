package sorts

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuickSort(t *testing.T) {
	a := []int{5, 4, 3, 2, 1}
	QuickSort(a, 0, len(a)-1)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}
