package rmq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumericBinaryIndexedTree_SumRange(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	var b = BuildSumBIT(a)
	assert.Equal(t, 7, b.SumRange(2, 3), "Sum failed")
}
