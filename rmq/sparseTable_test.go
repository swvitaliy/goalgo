package rmq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLg(t *testing.T) {
	assert.Equal(t, 0, lg(1))
	assert.Equal(t, 1, lg(2))
	assert.Equal(t, 2, lg(4))
	assert.Equal(t, 3, lg(8))
	assert.Equal(t, 3, lg(11))
	assert.Equal(t, 3, lg(15))
	assert.Equal(t, 4, lg(16))
	assert.Equal(t, 4, lg(17))
}

func TestSparseTable_Rmq(t *testing.T) {
	b := []int{0, 1, 3, 4, 5, 6, 7}
	st := BuildSparseTable(b, maxFn)
	assert.Equal(t, 1, st.Rmq(0, 1, maxFn))
	assert.Equal(t, 4, st.Rmq(0, 3, maxFn))
	assert.Equal(t, 4, st.Rmq(2, 3, maxFn))
	assert.Equal(t, 6, st.Rmq(3, 5, maxFn))
}

func maxFn(a, b int) int {
	return max(a, b)
}
