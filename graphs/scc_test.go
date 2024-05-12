package graphs

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestScc_Single(t *testing.T) {
	g := [][]int{
		{1, 2},
		{0, 2},
		{0, 1, 3},
		{2},
	}

	cs := map[int][]int{
		0: {0, 1, 2, 3},
	}

	used := make([]bool, 4)

	gr := MakeReversedGraph(g, 4)
	Scc(g, gr, 4, func(comp []int) {
		sort.Ints(comp)
		t.Log(comp)
		assert.False(t, used[comp[0]])
		used[comp[0]] = true
		assert.Equal(t, cs[comp[0]], comp)
	})
}

func TestScc_FourComps(t *testing.T) {
	g := [][]int{
		{1, 2},
		{0, 2},
		{0, 1},
		{4},
		{5},
		{3},
		{6, 7},
		{6},
		{},
	}

	cs := map[int][]int{
		0: {0, 1, 2},
		3: {3, 4, 5},
		6: {6, 7},
		8: {8},
	}

	used := make([]bool, len(g))

	gr := MakeReversedGraph(g, len(g))
	Scc(g, gr, len(g), func(comp []int) {
		sort.Ints(comp)
		t.Log(comp)
		assert.False(t, used[comp[0]])
		used[comp[0]] = true
		assert.Equal(t, cs[comp[0]], comp)
	})
}
