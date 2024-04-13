package rmq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumericFenwickTree_SumRange(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	var b = BuildSumFenwickTree(a)
	assert.Equal(t, 7, b.SumRange(2, 3), "Sum failed")
}

func TestNumericFenwickTree_SumRange2(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	var b = BuildFenwickTree(a, sumFn)
	assert.Equal(t, 6, b.Get(2, sumFn), "Sum failed")
	assert.Equal(t, 7, b.GetRange(2, 3, sumFn, subFn), "Sum failed")
}
