package consistent

import (
	"math"
	"strconv"
	"testing"
)

func TestHashRing_Lookup(t *testing.T) {
	// Test the key distribution of the consistent hash ring
	ring := NewHashRing(3)
	ring.Add("node1")
	ring.Add("node2")
	ring.Add("node3")

	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range keys {
		node := ring.Lookup([]byte(key))
		if node == "" {
			t.Errorf("Expected a node for key %s, but got none", key)
		}
	}
}

func TestHashRing_Distribution(t *testing.T) {
	ring := NewHashRing(100)
	ring.Add("node1")
	ring.Add("node2")
	ring.Add("node3")

	counts := make(map[string]int)
	totalCalls := 1000000
	for i := 0; i < totalCalls; i++ {
		key := []byte(strconv.Itoa(i))
		node := ring.Lookup(key)
		counts[node]++
	}

	expected := 1.0 / 3.0
	eps := 0.1
	for node, count := range counts {
		actual := float64(count) / float64(totalCalls)
		if math.Abs(actual-expected) > eps {
			t.Errorf("Node %s: expected ~%.3f, got %.3f", node, expected, actual)
		}
	}
}
