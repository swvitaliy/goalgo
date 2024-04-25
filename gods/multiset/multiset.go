package multiset

import (
	"github.com/emirpasic/gods/v2/trees/btree"
)

func MsPut[K comparable](t *btree.Tree[K, int], k K) {
	node := t.GetNode(k)
	if node == nil {
		t.Put(k, 1)
	} else {
		i, ok := search(t, node, k)
		if !ok {
			panic("multiset put failed - wrong node")
		}
		node.Entries[i].Value++
	}
}

func MsRemove[K comparable](t *btree.Tree[K, int], k K) {
	node := t.GetNode(k)
	if node == nil {
		panic("multiset remove failed - node not found")
	} else {
		i, ok := search(t, node, k)
		if !ok {
			panic("multiset remove failed - wrong node")
		}
		node.Entries[i].Value--
		if node.Entries[i].Value == 0 {
			t.Remove(k)
		}
	}
}

// search run binary search for entry on a node
// It has copied from btree sources (v2.0.0-alpha, 21.04.2024)
func search[K comparable, V any](tree *btree.Tree[K, V], node *btree.Node[K, V], key K) (index int, found bool) {
	low, high := 0, len(node.Entries)-1
	var mid int
	for low <= high {
		mid = (high + low) / 2
		compare := tree.Comparator(key, node.Entries[mid].Key)
		switch {
		case compare > 0:
			low = mid + 1
		case compare < 0:
			high = mid - 1
		case compare == 0:
			return mid, true
		}
	}
	return low, false
}
