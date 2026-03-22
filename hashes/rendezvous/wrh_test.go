package rendezvous

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// напиши новый тест, который будет вызывать Lookup 1_000_000 раз и считать для каждого  узла сколько он был получен.
// Затем, с точностью Eps сравнивал результаты с относительными весами

func TestWRH_Lookup_NoZeroWeighted(t *testing.T) {
	weightedNodes := map[string]float64{
		"A": 10,
		"B": 30,
		"C": 0,
	}
	wrh := NewWRH(weightedNodes)
	first := wrh.Lookup([]byte("first"))
	second := wrh.Lookup([]byte("second"))
	third := wrh.Lookup([]byte("third"))
	some := wrh.Lookup([]byte("1"))
	assert.Equal(t, first, "B")
	assert.Equal(t, second, "B")
	assert.Equal(t, third, "B")
	assert.Equal(t, some, "A")
}

func TestWRH_Distribution(t *testing.T) {
	weightedNodes := map[string]float64{
		"A": 10,
		"B": 30,
		"C": 0,
	}
	wrh := NewWRH(weightedNodes)
	counts := make(map[string]int)
	totalCalls := 1000000
	for i := 0; i < totalCalls; i++ {
		key := []byte(strconv.Itoa(i))
		node := wrh.Lookup(key)
		counts[node]++
	}
	totalWeight := 10.0 + 30.0 // 40
	expectedA := 10.0 / totalWeight
	expectedB := 30.0 / totalWeight
	actualA := float64(counts["A"]) / float64(totalCalls)
	actualB := float64(counts["B"]) / float64(totalCalls)
	eps := 0.01 // точность
	assert.InDelta(t, expectedA, actualA, eps)
	assert.InDelta(t, expectedB, actualB, eps)
	// Убедиться, что C не выбран
	assert.Equal(t, 0, counts["C"])
}
