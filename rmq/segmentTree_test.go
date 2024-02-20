package rmq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildSegTree(t *testing.T) {
	a := BuildSegTree([]int{1, 2, 3, 4, 5}, maxFn)
	assert.Equal(t, 3, a.Query(0, 2, maxFn))
}

func TestSegTree_Query(t *testing.T) {
	a := BuildSegTree([]int{1, 2, 3, 4, 5}, sumFn)
	assert.Equal(t, 9, a.Query(1, 3, sumFn))
}

func sumFn(a, b int) int {
	return a + b
}
