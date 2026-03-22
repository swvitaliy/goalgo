package rendezvous

import "github.com/cespare/xxhash/v2"

type HashFunc[K comparable] func(key K) uint64

// HRW implements the HRW Rendezvous algorithm
// for selecting a node based on a key.
// Inspired by https://github.com/dgryski/go-rendezvous/blob/master/rdv.go
type HRW struct {
	nhashes []uint64
}

// NewHRW creates a new HRW instance with the given nodes and hash function.
func NewHRW(nodes []string) *HRW {
	return &HRW{
		nhashes: sliceMap(nodes),
	}
}

func sliceMap(a []string) []uint64 {
	res := make([]uint64, len(a))
	for i, v := range a {
		res[i] = xxhash.Sum64String(v)
	}
	return res
}

// Lookup returns the index of the node that should be selected for the given key.
func (h *HRW) Lookup(key []byte) int {
	var maxHash uint64 = 0
	maxIndex := -1
	khash := xxhash.Sum64(key)
	for i, nhash := range h.nhashes {
		hash := xorshiftMult64(khash ^ nhash)
		if hash > maxHash {
			maxHash = hash
			maxIndex = i
		}
	}

	return maxIndex
}

func xorshiftMult64(x uint64) uint64 {
	x ^= x >> 12 // a
	x ^= x << 25 // b
	x ^= x >> 27 // c
	return x * 2685821657736338717
}
