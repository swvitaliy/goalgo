package consistent

import (
	"sort"
	"sync"

	"github.com/cespare/xxhash/v2"
)

const (
	prime = 16777619 // A prime number used for hashing
)

type HashRing struct {
	replicas int
	nhashes  []uint64
	ring     map[uint64]string
	mu       sync.RWMutex
}

// NewHashRing creates a new consistent hash instance with the given nodes and hash function.
func NewHashRing(replicas int) *HashRing {
	h := &HashRing{
		replicas: replicas,
		nhashes:  make([]uint64, 0),
		ring:     make(map[uint64]string),
	}
	return h
}

// Add adds a new node to the consistent hash ring.
func (h *HashRing) Add(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < h.replicas; i++ {
		hash := xorshiftMult64(xxhash.Sum64String(node) ^ (uint64(i) * prime))
		h.ring[hash] = node
		h.nhashes = append(h.nhashes, hash)
	}

	sort.Slice(h.nhashes, func(i, j int) bool {
		return h.nhashes[i] < h.nhashes[j]
	})
}

// Remove removes a node from the consistent hash ring.
func (h *HashRing) Remove(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < h.replicas; i++ {
		hash := xorshiftMult64(xxhash.Sum64String(node) ^ (uint64(i) * prime))
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
func (h *HashRing) Lookup(key []byte) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.nhashes) == 0 {
		return ""
	}

	hash := xxhash.Sum64(key)

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
