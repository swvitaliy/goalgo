package rendezvous

import (
	"math"

	"github.com/cespare/xxhash/v2"
)

type wrhNode struct {
	key       string
	hash      uint64
	invWeight float64
}

// WRH implements the Weighted Rendezvous Hashing algorithm
// for selecting a node based on a key.
// Inspired by https://github.com/dgryski/go-rendezvous/blob/master/rdv.go
type WRH struct {
	nodes []wrhNode
}

// NewWRH creates a new Weighted Rendezvous Hash instance with the given nodes and hash function.
func NewWRH(nodes map[string]float64) *WRH {
	return &WRH{
		nodes: sliceOfNodes(nodes),
	}
}

func sliceOfNodes(nodes map[string]float64) []wrhNode {
	result := make([]wrhNode, len(nodes))
	i := 0
	for k, v := range nodes {
		invWeight := 0.0
		if v > 0 {
			invWeight = 1.0 / v
		}
		result[i] = wrhNode{k, xxhash.Sum64String(k), invWeight}
		i++
	}
	return result
}

// Lookup returns the index of the node that should be selected for the given key.
func (h *WRH) Lookup(key []byte) string {
	var minScore float64 = math.Inf(1)
	var resultNode string
	khash := xxhash.Sum64(key)
	for _, node := range h.nodes {
		if node.invWeight == 0.0 {
			continue
		}
		hash := xorshiftMult64(khash ^ node.hash)
		u := float64(hash) / float64(math.MaxUint64)
		s := -math.Log(u) * node.invWeight
		if s < minScore {
			minScore = s
			resultNode = node.key
		}
	}

	return resultNode
}
