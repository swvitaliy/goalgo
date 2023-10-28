package goalgo

type lruNode struct {
	Key  int
	Val  int
	Next *lruNode
	Prev *lruNode
}

type LRUCache struct {
	Data       map[int]*lruNode
	Head, Tail *lruNode
	Cap, Cnt   int
}

func NewLRUCache(c int) LRUCache {
	lru := LRUCache{}
	lru.Cap = c
	lru.Cnt = 0
	lru.Data = make(map[int]*lruNode)
	lru.Head = &lruNode{}
	lru.Tail = &lruNode{}
	lru.Head.Next = lru.Tail
	lru.Tail.Prev = lru.Head
	return lru
}

func (lru *LRUCache) addNode(n *lruNode) {
	next := lru.Head.Next

	n.Next = next
	n.Prev = lru.Head

	next.Prev = n
	lru.Head.Next = n
}

func (lru *LRUCache) removeNode(n *lruNode) {
	prev := n.Prev
	next := n.Next

	prev.Next = next
	next.Prev = prev
}

func (lru *LRUCache) moveToHead(n *lruNode) {
	lru.removeNode(n)
	lru.addNode(n)
}

func (lru *LRUCache) popTail() *lruNode {
	t := lru.Tail.Prev
	lru.removeNode(t)
	return t
}

func (lru *LRUCache) Get(key int) int {
	n, ok := lru.Data[key]
	if !ok {
		return -1
	}

	lru.moveToHead(n)
	return n.Val
}

func (lru *LRUCache) Put(key int, value int) {
	n, ok := lru.Data[key]
	if ok {
		n.Val = value
		lru.moveToHead(n)
	} else {
		n := &lruNode{Key: key, Val: value}
		lru.Data[key] = n
		lru.addNode(n)
		lru.Cnt++

		if lru.Cnt > lru.Cap {
			t := lru.popTail()
			delete(lru.Data, t.Key)
			lru.Cnt--
		}
	}
}
