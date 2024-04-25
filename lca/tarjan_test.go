package lca

import "testing"

func TestGraphLca(t *testing.T) {
	g := [][]int{
		[]int{1, 2, 3},
		[]int{0, 4},
		[]int{0, 5},
		[]int{0, 6},
		[]int{2, 7},
		[]int{2, 8},
		[]int{3, 9},
		[]int{3, 10},
		[]int{7, 11},
		[]int{8, 12},
		[]int{9, 13},
	}

	q := [][]Query{
		[]Query{{from: 0, to: -1}},
		[]Query{{from: 1, to: -1}},
		[]Query{{from: 2, to: -1}},
		[]Query{{from: 3, to: -1}},
		[]Query{{from: 4, to: -1}},
		[]Query{{from: 5, to: -1}},
		[]Query{{from: 6, to: -1}},
		[]Query{{from: 7, to: -1}},
		[]Query{{from: 8, to: -1}},
		[]Query{{from: 9, to: -1}},
		[]Query{{from: 10, to: -1}},
		[]Query{{from: 11, to: -1}},
		[]Query{{from: 12, to: -1}},
		[]Query{{from: 13, to: -1}},
	}

	GraphLca(g, q)
}
