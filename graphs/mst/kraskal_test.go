package mst

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

// Returns a list of edges representing the minimum spanning tree of the input graph.
func TestKraskalMST_Cycle(t *testing.T) {
	g := []edge{
		{from: 0, to: 1, w: 5},
		{from: 1, to: 2, w: 3},
		{from: 2, to: 0, w: 1},
	}
	n := 3

	expected := []edge{
		{from: 2, to: 0, w: 1},
		{from: 1, to: 2, w: 3},
	}

	result := KraskalMST(g, n)

	sort.Slice(result, func(i, j int) bool {
		return result[i].w < result[j].w
	})

	if len(result) != len(expected) {
		t.Errorf("Expected %d edges, but got %d", len(expected), len(result))
	}

	for i := 0; i < len(result); i++ {
		if result[i].from != expected[i].from || result[i].to != expected[i].to || result[i].w != expected[i].w {
			t.Errorf("Expected edge %v, but got %v", expected[i], result[i])
		}
	}
}

// Returns an empty list when given an empty graph.
func TestKraskalMST_EmptyGraph(t *testing.T) {
	g := []edge{}
	n := 0

	expected := []edge{}

	result := KraskalMST(g, n)

	if len(result) != len(expected) {
		t.Errorf("Expected %d edges, but got %d", len(expected), len(result))
	}
}

// Returns a single edge when given a graph with two nodes and one edge.
func TestKraskalMST_SingleEdgeGraph(t *testing.T) {
	g := []edge{
		{from: 0, to: 1, w: 5},
	}
	n := 2

	expected := []edge{
		{from: 0, to: 1, w: 5},
	}

	result := KraskalMST(g, n)
	assert.Equal(t, expected, result)
}

// Returns an empty list when given a disconnected graph.
func TestKraskalMST_DisconnectedGraph(t *testing.T) {
	g := []edge{
		{from: 0, to: 1, w: 5},
		{from: 2, to: 3, w: 2},
	}
	n := 4

	expected := []edge{
		{from: 2, to: 3, w: 2},
		{from: 0, to: 1, w: 5},
	}

	result := KraskalMST(g, n)

	if len(result) != len(expected) {
		t.Errorf("Expected %d edges, but got %d", len(expected), len(result))
	}
}

// Returns an empty list when given a graph with one node.
func TestKraskalMST_SingleNodeGraph(t *testing.T) {
	g := []edge{}
	n := 1

	expected := []edge{}

	result := KraskalMST(g, n)

	if len(result) != len(expected) {
		t.Errorf("Expected %d edges, but got %d", len(expected), len(result))
	}
}

// Returns an empty list when given a graph with two nodes and no edges.
func TestKraskalMST_NoEdgeGraph(t *testing.T) {
	g := []edge{}
	n := 2

	expected := []edge{}

	result := KraskalMST(g, n)

	if len(result) != len(expected) {
		t.Errorf("Expected %d edges, but got %d", len(expected), len(result))
	}
}
