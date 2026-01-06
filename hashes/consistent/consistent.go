package consistent

import (
	"sort"
	"sync"
)

const (
	prime = 16777619 // A prime number used for hashing
)

type HashFunc[K comparable] func(key K) uint64

type HashRing[K comparable] struct {
	hash     HashFunc[K]
	replicas int
	nhashes  []uint64
	ring     map[uint64]K
	mu       sync.RWMutex
}

// NewHashingRing creates a new consistent hash instance with the given nodes and hash function.
func NewHashingRing[K comparable](nodes []K, replicas int, hash HashFunc[K]) *HashRing[K] {
	h := &HashRing[K]{
		hash:     hash,
		replicas: replicas,
		nhashes:  make([]uint64, 0, replicas*len(nodes)),
		ring:     make(map[uint64]K, len(nodes)),
	}

	for _, node := range nodes {
		h.Add(node)
	}

	return h
}

// Add adds a new node to the consistent hash ring.
func (h *HashRing[K]) Add(node K) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < h.replicas; i++ {
		hash := xorshiftMult64(h.hash(node) ^ (uint64(i) * prime))
		h.ring[hash] = node
		h.nhashes = append(h.nhashes, hash)
	}

	sort.Slice(h.nhashes, func(i, j int) bool {
		return h.nhashes[i] < h.nhashes[j]
	})
}

// Remove removes a node from the consistent hash ring.
func (h *HashRing[K]) Remove(node K) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < h.replicas; i++ {
		hash := xorshiftMult64(h.hash(node) ^ (uint64(i) * prime))
		delete(h.ring, hash)
		for j, nhash := range h.nhashes {
			if nhash == hash {
				h.nhashes = append(h.nhashes[:j], h.nhashes[j+1:]...)
				break
			}
		}
	}
}

// Lookup returns the node for the given key.
func (h *HashRing[K]) Lookup(key K) K {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.nhashes) == 0 {
		return *new(K) // Return zero value of K if no nodes are present
	}

	hash := h.hash(key)

	// Find the first hash that is greater than or equal to the key hash
	i := sort.Search(len(h.nhashes), func(i int) bool {
		return h.nhashes[i] >= hash
	})

	if i == len(h.nhashes) {
		i = 0
	}

	return h.ring[h.nhashes[i]]
}

func xorshiftMult64(x uint64) uint64 {
	x ^= x >> 12 // a
	x ^= x << 25 // b
	x ^= x >> 27 // c
	return x * 2685821657736338717
}
