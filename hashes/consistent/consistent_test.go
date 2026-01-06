package consistent

import (
	"fmt"
	"testing"
)

func KeyDestributionTest(t *testing.T) {
	// Test the key distribution of the consistent hash ring
	ring := NewHashRing(3, nil)
	ring.Add("node1")
	ring.Add("node2")
	ring.Add("node3")

	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range keys {
		node := ring.Get(key)
		if node == "" {
			t.Errorf("Expected a node for key %s, but got none", key)
		}
	}

	counts := map[string]int{}
	for i := 0; i < 100000; i++ {
		key := fmt.Sprintf("key%d", i)
		node := ch.Get(key)
		counts[node]++
	}
	fmt.Println(counts)
}
