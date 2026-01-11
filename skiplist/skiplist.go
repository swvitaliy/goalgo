package skiplist

import (
	"cmp"
	"iter"
	"runtime"
	"unsafe"
)

type levelGenerator interface {
	NextLevel() uint64
}

type Key cmp.Ordered
type Value any

type Node[K Key, V Value] struct {
	key   K
	value V
	next  []*Node[K, V]
}

type SkipList[K Key, V Value] struct {
	level          uint64
	size           uint64
	head           *Node[K, V]
	levelGenerator levelGenerator
}

func NewSkipList[K Key, V Value](lg levelGenerator) *SkipList[K, V] {
	return &SkipList[K, V]{
		level: 1,
		size:  0,
		head: &Node[K, V]{
			next: make([]*Node[K, V], maxLevel),
		},
		levelGenerator: lg,
	}
}

func NewFromSliceOfKeys[K Key](lg levelGenerator, a []K) *SkipList[K, struct{}] {
	sl := NewSkipList[K, struct{}](lg)
	for _, v := range a {
		sl.Insert(v, struct{}{})
	}
	return sl
}

func NewFromSliceOfValues[V Value](lg levelGenerator, a []V) *SkipList[int, V] {
	sl := NewSkipList[int, V](lg)
	for i, v := range a {
		sl.Insert(i, v)
	}
	return sl
}

type Pair[K Key, V Value] struct {
	Key   K
	Value V
}

func NewFromSliceOfPairs[K Key, V Value](lg levelGenerator, a []Pair[K, V]) *SkipList[K, V] {
	sl := NewSkipList[K, V](lg)
	for _, p := range a {
		sl.Insert(p.Key, p.Value)
	}
	return sl
}

func NewFromKeysIter[K Key](lg levelGenerator, it iter.Seq[K]) *SkipList[K, struct{}] {
	sl := NewSkipList[K, struct{}](lg)
	for v := range it {
		sl.Insert(v, struct{}{})
	}
	return sl
}

func NewFromValuesIter[V Value](lg levelGenerator, it iter.Seq[V]) *SkipList[int, V] {
	sl := NewSkipList[int, V](lg)
	i := 0
	for v := range it {
		sl.Insert(i, v)
		i++
	}
	return sl
}

func NewFromPairsIter[K Key, V Value](lg levelGenerator, it iter.Seq[Pair[K, V]]) *SkipList[K, V] {
	sl := NewSkipList[K, V](lg)
	for p := range it {
		sl.Insert(p.Key, p.Value)
	}
	return sl
}

// SearchNode ищет узел по ключу
func (sl *SkipList[K, V]) SearchNode(target K) *Node[K, V] {
	node := sl.head
	for i := sl.level - 1; i >= 0; i-- {
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

func (sl *SkipList[K, V]) Search(target K) (value V, ok bool) {
	node := sl.SearchNode(target)
	if node == nil {
		ok = false
		return
	}
	return node.value, true
}

// BatchSearchNodes ищет массив ключей с prefetch
func (sl *SkipList[K, V]) BatchSearchNodes(keys []K) []*Node[K, V] {
	prefetch := func(n *Node[K, V]) {
		if n != nil {
			runtime.KeepAlive(n)
			_ = *(*uintptr)(unsafe.Pointer(n))
		}
	}

	results := make([]*Node[K, V], len(keys))
	for i, key := range keys {
		curr := sl.head
		for lvl := sl.level - 1; lvl >= 0; lvl-- {
			next := curr.next[lvl]
			if lvl == 0 {
				prefetch(next)
			}
			for next != nil && next.key < key {
				curr = next
				next = curr.next[lvl]
				if lvl == 0 {
					prefetch(next)
				}
			}
		}
		results[i] = curr.next[0]
	}
	return results
}

func (sl *SkipList[K, V]) Insert(key K, value V) *Node[K, V] {
	update := make([]*Node[K, V], maxLevel)
	node := sl.head
	for i := sl.level - 1; i >= 0; i-- {
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
	lvl := sl.levelGenerator.NextLevel()
	if lvl > sl.level {
		for i := sl.level; i < lvl; i++ {
			update[i] = sl.head
		}
		sl.level = lvl
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
	sl.size++
	return newNode
}

func (sl *SkipList[K, V]) Delete(key K) bool {
	update := make([]*Node[K, V], maxLevel)
	node := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for node.next[i] != nil && node.next[i].key < key {
			node = node.next[i]
		}
		update[i] = node
	}

	node = node.next[0]
	if node == nil || node.key != key {
		return false
	}

	for i := uint64(0); i < sl.level; i++ {
		if update[i].next[i] != node {
			break
		}
		update[i].next[i] = node.next[i]
	}

	for sl.level > 1 && sl.head.next[sl.level-1] == nil {
		sl.level--
	}
	sl.size--
	return true
}

func (sl *SkipList[K, V]) Range(b, e K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		node := sl.head
		for i := sl.level - 1; i >= 0; i-- {
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
