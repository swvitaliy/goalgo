package skiplist

import (
	"iter"
	"sort"
	"sync/atomic"
	"unsafe"
)

type ConcurrentNode[K Key, V Value] struct {
	key   K
	value V
	next  []atomic.Pointer[ConcurrentNode[K, V]]
	level int32
}

type ConcurrentSkipList[K Key, V Value] struct {
	head           *ConcurrentNode[K, V]
	level          int32
	levelGenerator levelGenerator
}

func NewConcurrentSkipList[K Key, V Value](lg levelGenerator) *ConcurrentSkipList[K, V] {
	var zeroKey K
	var zeroVal V
	head := newNode[K, V](zeroKey, zeroVal, maxLevel)
	return &ConcurrentSkipList[K, V]{
		head:           head,
		levelGenerator: lg,
	}
}

func NewConcurrentFromKeysIter[K Key](lg levelGenerator, it iter.Seq[K]) *ConcurrentSkipList[K, struct{}] {
	sl := NewConcurrentSkipList[K, struct{}](lg)
	for v := range it {
		sl.Insert(v, struct{}{})
	}
	return sl
}

func NewConcurrentFromValuesIter[V Value](lg levelGenerator, it iter.Seq[V]) *ConcurrentSkipList[int, V] {
	sl := NewConcurrentSkipList[int, V](lg)
	i := 0
	for v := range it {
		sl.Insert(i, v)
		i++
	}
	return sl
}

func NewConcurrentFromPairsIter[K Key, V Value](lg levelGenerator, it iter.Seq[Pair[K, V]]) *ConcurrentSkipList[K, V] {
	sl := NewConcurrentSkipList[K, V](lg)
	for p := range it {
		sl.Insert(p.Key, p.Value)
	}
	return sl
}

func newNode[K Key, V Value](key K, value V, level int) *ConcurrentNode[K, V] {
	n := &ConcurrentNode[K, V]{
		key:   key,
		value: value,
		level: int32(level),
		next:  make([]atomic.Pointer[ConcurrentNode[K, V]], level),
	}
	for i := range n.next {
		n.next[i].Store(nil)
	}
	return n
}

// find locates preds and succs for key. Also helps unlink marked nodes.
// returns true if key is found (and succs[0] has key == key and not marked)
func (csl *ConcurrentSkipList[K, V]) find(key K, preds, succs []*ConcurrentNode[K, V]) bool {
	type node = ConcurrentNode[K, V]

	var (
		pred *node
		curr *node
		succ *node
		mark bool
	)
retry:
	pred = csl.head
	for level := int(maxLevel - 1); level >= 0; level-- {
		curr, _ = loadNext(pred, level)
		for {
			if curr == nil {
				break
			}
			succ, mark = loadNext(curr, level)
			for mark {
				// attempt physical removal
				if !casNext(pred, int32(level), curr, false, succ, false) {
					// failed - start over
					goto retry
				}
				curr = succ
				if curr == nil {
					break
				}
				succ, mark = loadNext(curr, level)
			}
			if curr == nil {
				break
			}
			if curr.key < key {
				pred = curr
				curr = succ
			} else {
				break
			}
		}
		preds[level] = pred
		succs[level] = curr
	}
	// check if found at level 0 and not marked
	if curr != nil && curr.key == key {
		_, mark = loadNext(curr, 0)
		return !mark
	}
	return false
}

// Contains / Search
func (csl *ConcurrentSkipList[K, V]) Contains(key K) (*V, bool) {
	type node = ConcurrentNode[K, V]

	preds := make([]*node, maxLevel)
	succs := make([]*node, maxLevel)
	found := csl.find(key, preds, succs)
	if !found {
		return nil, false
	}
	n := succs[0]
	if n == nil {
		return nil, false
	}
	return &n.value, true
}

// Insert (returns true if inserted, false if key already present)
func (csl *ConcurrentSkipList[K, V]) Insert(key K, value V) bool {
	type node = ConcurrentNode[K, V]
	level := int(csl.levelGenerator.NextLevel())
	var preds = make([]*node, maxLevel)
	var succs = make([]*node, maxLevel)

	for {
		found := csl.find(key, preds, succs)
		if found {
			// already present
			return false
		}
		newNode := newNode(key, value, level)
		for i := 0; i <= level; i++ {
			newNode.next[i].Store((*ConcurrentNode[K, V])(packPointer(succs[i], false)))
		}
		// try link at level 0 first
		pred := preds[0]
		succ := succs[0]
		if !casNext(pred, 0, succ, false, newNode, false) {
			// failed, retry
			continue
		}
		// link higher levels
		for i := 1; i <= level; i++ {
			for {
				pred = preds[i]
				succ = succs[i]
				if casNext(pred, int32(i), succ, false, newNode, false) {
					break
				}
				// if fail, recompute preds/succs
				csl.find(key, preds, succs)
			}
		}
		// Update list-level hint
		currentLevel := csl.level
		if int32(level) > currentLevel {
			atomic.CompareAndSwapInt32(&csl.level, currentLevel, int32(level))
		}
		return true
	}
}

func (csl *ConcurrentSkipList[K, V]) BulkInsert(keys []K, values []V) {
	if len(keys) != len(values) {
		panic("keys and values length mismatch")
	}

	type insertNode struct {
		key   K
		value V
		level int
	}

	nodes := make([]insertNode, len(keys))
	for i := range keys {
		nodes[i] = insertNode{
			key:   keys[i],
			value: values[i],
			level: int(csl.levelGenerator.NextLevel()),
		}
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].key < nodes[j].key })

	type nodeType = ConcurrentNode[K, V]

	// Shared traversal
	update := make([]*nodeType, maxLevel)
	current := csl.head

	for _, n := range nodes {
		for i := maxLevel - 1; i >= 0; i-- {
			for next := current.next[i].Load(); next != nil && next.key < n.key; next = current.next[i].Load() {
				current = next
			}
			update[i] = current
		}

		node := &ConcurrentNode[K, V]{
			key:   n.key,
			value: n.value,
			next:  make([]atomic.Pointer[nodeType], n.level),
		}

		for i := 0; i < n.level; i++ {
			node.next[i] = atomic.Pointer[nodeType]{}
			for {
				next := update[i].next[i].Load()
				node.next[i].Store(next)
				if update[i].next[i].CompareAndSwap(next, node) {
					break
				}
			}
		}
		current = node
	}
}

// Delete (returns true if deleted)
func (csl *ConcurrentSkipList[K, V]) Delete(key K) bool {
	type node = ConcurrentNode[K, V]
	var preds = make([]*node, maxLevel)
	var succs = make([]*node, maxLevel)
	var nodeToDelete *node
	for {
		found := csl.find(key, preds, succs)
		if !found {
			return false
		}
		nodeToDelete = succs[0]
		// logically mark from top level down to 0
		for level := nodeToDelete.level; level >= 1; level-- {
			var succ *node
			for {
				succ, _ = loadNext(nodeToDelete, int(level))
				if succ == nil {
					break
				}
				// mark pointer at this level
				if casNext(nodeToDelete, level, succ, false, succ, true) {
					break
				}
				// else retry
			}
		}
		// finally mark level 0
		for {
			succ, marked := loadNext(nodeToDelete, 0)
			if marked {
				// already marked by another remover
				return false
			}
			if casNext(nodeToDelete, 0, succ, false, succ, true) {
				// successful logical deletion
				break
			}
		}
		// try to physically remove by swinging preds' pointers
		csl.find(key, preds, succs) // helps unlink
		return true
	}
}

func (csl *ConcurrentSkipList[K, V]) BulkDelete(keys []K) {
	if len(keys) == 0 {
		return
	}

	// сортируем ключи для shared traversal
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	update := make([]*ConcurrentNode[K, V], maxLevel)
	current := csl.head

	for _, key := range keys {
		// shared traversal по уровням
		for i := int(maxLevel - 1); i >= 0; i-- {
			for {
				next, marked := loadNext(current, i)
				if next == nil || next.key >= key {
					break
				}
				current = next
				if marked {
					// пропускаем logically deleted узлы
					continue
				}
			}
			update[i] = current
		}

		// target node на уровне 0
		target, _ := loadNext(current, 0)
		if target == nil || target.key != key {
			// ключ не найден
			continue
		}

		// помечаем узел как удалённый на уровне 0
		for {
			ok := casNext(current, 0, target, false, target, true)
			if ok {
				break
			}
			// если кто-то уже пометил узел, выходим
			_, alreadyMarked := loadNext(current, 0)
			if alreadyMarked {
				break
			}
		}

		// unlink верхних уровней
		for i := int64(len(target.next) - 1); i >= 1; i-- { // верхние уровни
			for {
				next, marked := loadNext(target, int(i))
				if marked {
					next = nil
				}
				if casNext(update[i], int32(i), target, false, next, false) {
					break
				}
			}
		}

		// продолжаем traversal с последнего удалённого узла
		current = target
	}
}

// Compact - physically removes all marked nodes from the skip list
func (csl *ConcurrentSkipList[K, V]) Compact() {
	for level := int(atomic.LoadInt32(&csl.level)) - 1; level >= 0; level-- {
		prev := csl.head
		curr := prev.next[level].Load()

		for curr != nil {
			next := curr.next[level].Load()

			cPtr, m := unpackPointer[ConcurrentNode[K, V]](unsafe.Pointer(next))
			if m {
				if prev.next[level].CompareAndSwap(cPtr, next) {
					curr = next
					continue
				} else {
					curr = prev.next[level].Load()
					continue
				}
			}

			prev = curr
			curr = next
		}
	}
}

func packPointer[T any](p *T, marked bool) unsafe.Pointer {
	up := uintptr(unsafe.Pointer(p))
	if marked {
		up |= 1
	} else {
		up &^= 1
	}
	return unsafe.Pointer(up)
}

func unpackPointer[T any](p unsafe.Pointer) (*T, bool) {
	up := uintptr(p)
	marked := (up & 1) == 1
	ptr := unsafe.Pointer(up &^ 1) // очистка бита
	return (*T)(ptr), marked
}

func loadNext[K Key, V Value](n *ConcurrentNode[K, V], level int) (*ConcurrentNode[K, V], bool) {
	p := n.next[level].Load()
	if p == nil {
		return nil, false
	}
	return unpackPointer[ConcurrentNode[K, V]](unsafe.Pointer(p))
}

func casNext[K Key, V Value](n *ConcurrentNode[K, V], level int32, oldNode *ConcurrentNode[K, V], oldMarked bool, newNode *ConcurrentNode[K, V], newMarked bool) bool {
	type node = ConcurrentNode[K, V]
	oldPtr := packPointer(oldNode, oldMarked)
	newPtr := packPointer(newNode, newMarked)
	return n.next[level].CompareAndSwap((*node)(oldPtr), (*node)(newPtr))
}
