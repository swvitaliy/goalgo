package rendezvous

import "math"

type hashedNode[K comparable] struct {
	key    K
	nhash  uint64
	weight float64
}

// WRH implements the Weighted Rendezvous Hashing algorithm
// for selecting a node based on a key.
// Inspired by https://github.com/dgryski/go-rendezvous/blob/master/rdv.go
type WRH[K comparable] struct {
	hash    HashFunc[K]
	nhashes []hashedNode[K]
}

// NewWRH creates a new Weighted Rendezvous Hash instance with the given nodes and hash function.
func NewWRH[K comparable](nodes map[K]float64, hash HashFunc[K]) *WRH[K] {
	return &WRH[K]{
		hash:    hash,
		nhashes: convHashedNodes(nodes, hash),
	}
}

func convHashedNodes[K comparable](nodes map[K]float64, hash HashFunc[K]) []hashedNode[K] {
	nhashes := make([]hashedNode[K], len(nodes))
	for k, v := range nodes {
		nhashes[hash(k)] = hashedNode[K]{k, hash(k), v}
	}
	return nhashes
}

// Lookup returns the index of the node that should be selected for the given key.
func (h *WRH[K]) Lookup(key K) K {
	var maxScore float64 = 0
	var resultNode K
	khash := h.hash(key)
	for _, node := range h.nhashes {
		hash := xorshiftMult64(khash ^ node.nhash)
		u := float64(hash) / float64(math.MaxUint64)
		score := -math.Log(u) / node.weight
		if score > maxScore {
			maxScore = score
			resultNode = node.key
		}
	}

	return resultNode
}
