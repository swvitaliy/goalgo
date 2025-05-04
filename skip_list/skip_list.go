package skip_list

import (
	"cmp"
	"math/rand/v2"
)

const (
	p        = 0.5
	maxLevel = 16
)

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

func NewSkipList[K Key, V Value](maxLevel uint64) *SkipList[K, V] {
	return &SkipList[K, V]{
		level: 1,
		size:  0,
		head: &Node[K, V]{
			next: make([]*Node[K, V], maxLevel),
		},
	}
}

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
	lvl := randomLevel()
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

func randomLevel() uint64 {
	var lvl uint64 = 1
	for rand.Float64() < 0.5 && lvl < maxLevel {
		lvl++
	}
	return lvl
}
