package mst

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// returns a list of length n with the parent of each vertex in the MST
func TestPrimMST_ParentOfEachVertex(t *testing.T) {
	// Initialize the graph
	g := [][]adjEdge{
		{{t: 1, w: 2}, {t: 2, w: 3}},
		{{t: 0, w: 2}, {t: 2, w: 4}},
		{{t: 0, w: 3}, {t: 1, w: 4}},
	}

	// Invoke the PrimMST function
	result := PrimMST(0, g)

	// Check if the result has the correct length
	if len(result) != len(g) {
		t.Errorf("Expected result length %d, but got %d", len(g), len(result))
	}

	// Check if each vertex has a parent in the MST
	for i, parent := range result {
		if i == 0 && parent != -1 {
			t.Errorf("Expected parent of root vertex to be -1, but got %d", parent)
		} else if i != 0 && parent == -1 {
			t.Errorf("Expected parent of vertex %d to be non-negative, but got -1", i)
		}
	}
}

// returns a list of length n with -1 as parent for the root vertex
func TestPrimMST_RootVertexParent(t *testing.T) {
	// Initialize the graph
	g := [][]adjEdge{
		{{t: 1, w: 2}, {t: 2, w: 3}},
		{{t: 0, w: 2}, {t: 2, w: 4}},
		{{t: 0, w: 3}, {t: 1, w: 4}},
	}

	// Invoke the PrimMST function
	result := PrimMST(0, g)

	// Check if the parent of the root vertex is -1
	if result[0] != -1 {
		t.Errorf("Expected parent of root vertex to be -1, but got %d", result[0])
	}
}

// returns a list of length n with the root vertex as parent for a graph with only one vertex
func TestPrimMST_OneVertexGraph(t *testing.T) {
	// Initialize the graph
	g := [][]adjEdge{
		{},
	}

	// Invoke the PrimMST function
	result := PrimMST(0, g)

	assert.Equal(t, []ndx{-1}, result)
}

// returns an empty list for an empty graph
func TestPrimMST_EmptyGraph(t *testing.T) {
	// Initialize the graph
	g := [][]adjEdge{}

	// Invoke the PrimMST function
	result := PrimMST(0, g)

	assert.Equal(t, []ndx{}, result)
}
