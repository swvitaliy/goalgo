package rendezvous

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHRW(t *testing.T) {
	nodes := []string{"A", "B", "C"}
	weights := map[string]map[string]int{
		"A": {"X": 10, "Y": 20},
		"B": {"X": 30, "Y": 40},
		"C": {"X": 50, "Y": 60},
	}

	expected := map[string]string{
		"X": "C",
		"Y": "C",
	}

	result := NewHRW(nodes, weights)

	assert.Equal(t, expected, result)
}
