package rendezvous

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHRW(t *testing.T) {
	nodes := []string{"A", "B", "C"}
	hrw := NewHRW(nodes)

	first := hrw.Lookup([]byte("first"))
	second := hrw.Lookup([]byte("second"))
	third := hrw.Lookup([]byte("third"))
	assert.Equal(t, first, 1)
	assert.Equal(t, second, 1)
	assert.Equal(t, third, 0)
}

func TestHRW_Distribution(t *testing.T) {
	nodes := []string{"A", "B", "C"}
	hrw := NewHRW(nodes)

	counts := make(map[int]int)
	totalCalls := 1000000
	for i := 0; i < totalCalls; i++ {
		key := []byte(strconv.Itoa(i))
		nodeIndex := hrw.Lookup(key)
		counts[nodeIndex]++
	}

	expected := 1.0 / 3.0
	eps := 0.01
	for nodeIndex, count := range counts {
		actual := float64(count) / float64(totalCalls)
		assert.InDelta(t, expected, actual, eps, "Node index %d: expected ~%.3f, got %.3f", nodeIndex, expected, actual)
	}
}
