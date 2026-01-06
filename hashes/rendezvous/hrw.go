package rendezvous

type HashFunc[K comparable] func(key K) uint64

// HRW implements the HRW Rendezvous algorithm
// for selecting a node based on a key.
// Inspired by https://github.com/dgryski/go-rendezvous/blob/master/rdv.go
type HRW[K comparable] struct {
	hash    HashFunc[K]
	nhashes []uint64
}

// NewHRW creates a new HRW instance with the given nodes and hash function.
func NewHRW[K comparable](nodes []K, hash HashFunc[K]) *HRW[K] {
	return &HRW[K]{
		hash:    hash,
		nhashes: sliceMap(nodes, hash),
	}
}

func sliceMap[T any, R any](a []T, fn func(T) R) []R {
	res := make([]R, len(a))
	for i, v := range a {
		res[i] = fn(v)
	}
	return res
}

// Lookup returns the index of the node that should be selected for the given key.
func (h *HRW[K]) Lookup(key K) int {
	var maxHash uint64 = 0
	maxIndex := -1
	khash := h.hash(key)
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
