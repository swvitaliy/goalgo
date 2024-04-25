package lca

import (
	"github.com/stretchr/testify/assert"
	"goalgo/trees/bst"
	"testing"
)

func TestBstLca(t *testing.T) {
	root := bst.MakeFromLinear[int]([]int{6, 2, 8, 0, 4, 7, 9, -1, -1, 3, 5}, -1)
	lca := BstLca(root, 5, 9)
	assert.NotNil(t, lca)
	assert.Equal(t, 6, lca.Key)
}

func TestBstLcaNilRoot(t *testing.T) {
	lca := BstLca(nil, 5, 9)
	assert.Nil(t, lca)
}

func TestBstLcaOutOfRange(t *testing.T) {
	root := bst.MakeFromLinear[int]([]int{6, 2, 8, 0, 4, 7, 9, -1, -1, 3, 5}, -1)
	lca := BstLca(root, 10, 11)
	assert.Nil(t, lca)
}

func TestBstLcaNil(t *testing.T) {
	root := bst.MakeFromLinear[int]([]int{6, 2, 8, 0, 4, 7, 9, -1, -1, 3, 5}, -1)
	lca := BstLca(root, 10, 11)
	assert.Nil(t, lca)
}
