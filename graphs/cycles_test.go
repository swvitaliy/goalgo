package graphs

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestDetectCycles(t *testing.T) {
	g := [][]int{
		{1},
		{2},
		{0},
	}
	assert.True(t, DetectCycleDFS(g))
	cs := CyclesDFS(g)
	assert.Equal(t, 3, len(cs))
	sort.Ints(cs)
	assert.Equal(t, []int{0, 1, 2}, cs)
}

func TestNoCycle(t *testing.T) {
	g := [][]int{
		{1},
		{2},
		{},
	}
	assert.False(t, DetectCycleDFS(g))
}
