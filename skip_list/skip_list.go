package skip_list

import (
	"cmp"
	"iter"
	"math/rand/v2"
)

const (
	probability = 0.5
	maxLevel    = 16
)

func randomLevel() uint64 {
	var lvl uint64 = 1
	for rand.Float64() < probability && lvl < maxLevel {
		lvl++
	}
	return lvl
}

type Key cmp.Ordered
type Value any

type Node[K Key, V Value] struct {
	key   K
	value V
	next  []*Node[K, V]
}

type SkipList[K Key, V Value] struct {
	level uint64
	size  uint64
	head  *Node[K, V]
}

func NewSkipList[K Key, V Value]() *SkipList[K, V] {
	return &SkipList[K, V]{
		level: 1,
		size:  0,
		head: &Node[K, V]{
			next: make([]*Node[K, V], maxLevel),
		},
	}
}

func FromSliceOfKeys[K Key](a []K) *SkipList[K, struct{}] {
	sl := NewSkipList[K, struct{}]()
	for _, v := range a {
		sl.Insert(v, struct{}{})
	}
	return sl
}

func FromSliceOfValues[V Value](a []V) *SkipList[int, V] {
	sl := NewSkipList[int, V]()
	for i, v := range a {
		sl.Insert(i, v)
	}
	return sl
}

type Pair[K Key, V Value] struct {
	Key   K
	Value V
}

func FromSliceOfPairs[K Key, V Value](a []Pair[K, V]) *SkipList[K, V] {
	sl := NewSkipList[K, V]()
	for _, p := range a {
		sl.Insert(p.Key, p.Value)
	}
	return sl
}

func (csl *SkipList[K, V]) SearchNode(target K) *Node[K, V] {
	node := csl.head
	for i := csl.level - 1; i >= 0; i-- {
		for node != nil && node.next[i].key < target {
			node = node.next[i]
		}
	}
	if node == nil {
		return nil
	}
	node = node.next[0]
	if node != nil && node.key == target {
		return node
	}
	return nil
}

func (csl *SkipList[K, V]) Search(target K) (value V, ok bool) {
	node := csl.SearchNode(target)
	if node == nil {
		ok = false
		return
	}
	return node.value, true
}

func (csl *SkipList[K, V]) Insert(key K, value V) *Node[K, V] {
	update := make([]*Node[K, V], maxLevel)
	node := csl.head
	for i := csl.level - 1; i >= 0; i-- {
		for node != nil && node.next[i].key < key {
			node = node.next[i]
		}
		update[i] = node
	}

	if node == nil {
		return nil
	}

	node = node.next[0]
	if node != nil && node.key == key {
		node.value = value
		return node
	}
	lvl := randomLevel()
	if lvl > csl.level {
		for i := csl.level; i < lvl; i++ {
			update[i] = csl.head
		}
		csl.level = lvl
	}

	newNode := &Node[K, V]{
		key:   key,
		value: value,
		next:  make([]*Node[K, V], lvl),
	}
	for i := uint64(0); i < lvl; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
	csl.size++
	return newNode
}

func (csl *SkipList[K, V]) Delete(key K) bool {
	update := make([]*Node[K, V], maxLevel)
	node := csl.head
	for i := csl.level - 1; i >= 0; i-- {
		for node.next[i] != nil && node.next[i].key < key {
			node = node.next[i]
		}
		update[i] = node
	}

	node = node.next[0]
	if node == nil || node.key != key {
		return false
	}

	for i := uint64(0); i < csl.level; i++ {
		if update[i].next[i] != node {
			break
		}
		update[i].next[i] = node.next[i]
	}

	for csl.level > 1 && csl.head.next[csl.level-1] == nil {
		csl.level--
	}
	csl.size--
	return true
}

func (csl *SkipList[K, V]) Range(b, e K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		node := csl.head
		for i := csl.level - 1; i >= 0; i-- {
			for node.next[i] != nil && node.next[i].key < b {
				node = node.next[i]
			}
		}

		node = node.next[0]
		for node != nil && node.key <= e {
			if !yield(node.key, node.value) {
				return
			}
			node = node.next[0]
		}
	}
}
