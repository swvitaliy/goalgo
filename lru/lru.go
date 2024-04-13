package lru

type node[K comparable, V any] struct {
	Key  K
	Val  V
	Next *node[K, V]
	Prev *node[K, V]
}

type Cache[K comparable, V any] struct {
	data       map[K]*node[K, V]
	head, tail *node[K, V]
	cap, cnt   int
}

func New[K comparable, V any](c int) Cache[K, V] {
	lru := Cache[K, V]{}
	lru.cap = c
	lru.cnt = 0
	lru.data = make(map[K]*node[K, V])
	lru.head = &node[K, V]{}
	lru.tail = &node[K, V]{}
	lru.head.Next = lru.tail
	lru.tail.Prev = lru.head
	return lru
}

func (lru *Cache[K, V]) addNode(n *node[K, V]) {
	next := lru.head.Next

	n.Next = next
	n.Prev = lru.head

	next.Prev = n
	lru.head.Next = n
}

func (lru *Cache[K, V]) removeNode(n *node[K, V]) {
	prev := n.Prev
	next := n.Next

	prev.Next = next
	next.Prev = prev
}

func (lru *Cache[K, V]) moveToHead(n *node[K, V]) {
	lru.removeNode(n)
	lru.addNode(n)
}

func (lru *Cache[K, V]) popTail() *node[K, V] {
	t := lru.tail.Prev
	lru.removeNode(t)
	return t
}

func (lru *Cache[K, V]) Get(key K) (V, bool) {
	n, ok := lru.data[key]
	if !ok {
		var v V
		return v, false
	}

	lru.moveToHead(n)
	return n.Val, true
}

func (lru *Cache[K, V]) Put(key K, value V) {
	n, ok := lru.data[key]
	if ok {
		n.Val = value
		lru.moveToHead(n)
	} else {
		n := &node[K, V]{Key: key, Val: value}
		lru.data[key] = n
		lru.addNode(n)
		lru.cnt++

		if lru.cnt > lru.cap {
			t := lru.popTail()
			delete(lru.data, t.Key)
			lru.cnt--
		}
	}
}
