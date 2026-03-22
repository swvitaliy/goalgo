package hashes

import (
	"fmt"
	"strconv"
	"testing"

	"goalgo/hashes/consistent"
	"goalgo/hashes/rendezvous"
)

// benchmark of change distributions for different hash functions (hrw, wrh, consistent)

func BenchmarkHRW_Distribution(b *testing.B) {
	numNodesList := []int{3, 5, 10, 20, 50, 100}
	for _, numNodes := range numNodesList {
		b.Run(fmt.Sprintf("nodes_%d", numNodes), func(b *testing.B) {
			nodes := make([]string, numNodes)
			for i := 0; i < numNodes; i++ {
				nodes[i] = fmt.Sprintf("node%d", i)
			}
			hrw := rendezvous.NewHRW(nodes)
			counts := make(map[int]int)
			totalCalls := 10000000
			for i := 0; i < totalCalls; i++ {
				key := []byte(strconv.Itoa(i))
				nodeIndex := hrw.Lookup(key)
				counts[nodeIndex]++
			}
			expected := 1.0 / float64(numNodes)
			variance := 0.0
			for _, count := range counts {
				actual := float64(count) / float64(totalCalls)
				diff := actual - expected
				variance += diff * diff
			}
			variance /= float64(numNodes)
			b.Logf("HRW with %d nodes: variance from uniform = %.6f", numNodes, variance)
		})
	}
}

func BenchmarkWRH_Distribution(b *testing.B) {
	numNodesList := []int{3, 5, 10, 20, 50, 100}
	for _, numNodes := range numNodesList {
		b.Run(fmt.Sprintf("nodes_%d", numNodes), func(b *testing.B) {
			weightedNodes := make(map[string]float64)
			for i := 0; i < numNodes; i++ {
				weightedNodes[fmt.Sprintf("node%d", i)] = 1.0
			}
			wrh := rendezvous.NewWRH(weightedNodes)
			counts := make(map[string]int)
			totalCalls := 10000000
			for i := 0; i < totalCalls; i++ {
				key := []byte(strconv.Itoa(i))
				node := wrh.Lookup(key)
				counts[node]++
			}
			expected := 1.0 / float64(numNodes)
			variance := 0.0
			for _, count := range counts {
				actual := float64(count) / float64(totalCalls)
				diff := actual - expected
				variance += diff * diff
			}
			variance /= float64(numNodes)
			b.Logf("WRH with %d nodes: variance from uniform = %.6f", numNodes, variance)
		})
	}
}

func BenchmarkConsistent_Distribution(b *testing.B) {
	numNodesList := []int{3, 5, 10, 20, 50, 100}
	for _, numNodes := range numNodesList {
		b.Run(fmt.Sprintf("nodes_%d", numNodes), func(b *testing.B) {
			ring := consistent.NewHashRing(1) // 1 replica per node for uniformity
			for i := 0; i < numNodes; i++ {
				ring.Add(fmt.Sprintf("node%d", i))
			}
			counts := make(map[string]int)
			totalCalls := 10000000
			for i := 0; i < totalCalls; i++ {
				key := []byte(strconv.Itoa(i))
				node := ring.Lookup(key)
				counts[node]++
			}
			expected := 1.0 / float64(numNodes)
			variance := 0.0
			for _, count := range counts {
				actual := float64(count) / float64(totalCalls)
				diff := actual - expected
				variance += diff * diff
			}
			variance /= float64(numNodes)
			b.Logf("Consistent with %d nodes: variance from uniform = %.6f", numNodes, variance)
		})
	}
}
