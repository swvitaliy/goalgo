package lfu

import "container/list"

// two priorities: frequency and recently used.
// more frequency and more recently used are kept.
// linked list to store the key sequence with
// a map to store the head of different frequency

type Cache[Key comparable, Val any] struct {
	cap        int
	list       *list.List
	defaultVal Val
	stores     map[Key]*list.Element
	// key is the frequency, val is the element head (most recently used)
	frequencyHead map[int]*list.Element
}

type Item[Key comparable, Val any] struct {
	key       Key
	val       Val
	frequency int
}

func New[Key comparable, Val any](capacity int, defaultVal Val) Cache[Key, Val] {
	return Cache[Key, Val]{
		cap:           capacity,
		list:          list.New(),
		defaultVal:    defaultVal,
		stores:        make(map[Key]*list.Element, capacity),
		frequencyHead: make(map[int]*list.Element, capacity),
	}
}

func (cache *Cache[Key, Val]) Get(key Key) Val {
	val, ok := cache.stores[key]
	if !ok {
		return cache.defaultVal
	}
	item, _ := val.Value.(*Item[Key, Val])
	cache.refreshFrequency(cache.stores[key])
	return item.val
}

func (cache *Cache[Key, Val]) Put(key Key, value Val) {
	_, ok := cache.stores[key]
	if !ok {
		cache.clean()
		item := &Item[Key, Val]{key, value, 0}
		cache.stores[key] = cache.list.PushBack(item)
	} else {
		item := cache.stores[key].Value.(*Item[Key, Val])
		item.val = value
	}
	cache.refreshFrequency(cache.stores[key])
}

func (cache *Cache[Key, Val]) refreshFrequency(elem *list.Element) {
	item, _ := elem.Value.(*Item[Key, Val])
	oldHead, oldOk := cache.frequencyHead[item.frequency]
	// if not ok, it is the new element, since no key's frequency is 0, we do not worry
	// if ok and old head is the elem, we need to set old head to the next value if freqency is the same
	if oldOk && oldHead == elem {
		if oldHead.Next() != nil && oldHead.Next().Value.(*Item[Key, Val]).frequency == item.frequency {
			cache.frequencyHead[item.frequency] = oldHead.Next()
		} else {
			delete(cache.frequencyHead, item.frequency)
		}
	}

	item.frequency++

	newHead, newOk := cache.frequencyHead[item.frequency]
	if newOk {
		cache.list.MoveBefore(elem, newHead)
	} else if oldOk {
		cache.list.MoveBefore(elem, oldHead)
	}
	cache.frequencyHead[item.frequency] = elem
}

func (cache *Cache[Key, Val]) clean() {
	for cache.list.Len() >= cache.cap {
		back := cache.list.Back()
		item, _ := back.Value.(*Item[Key, Val])
		delete(cache.stores, item.key)
		if cache.frequencyHead[item.frequency] == back {
			delete(cache.frequencyHead, item.frequency)
		}
		cache.list.Remove(cache.list.Back())
	}
}

func (cache *Cache[Key, Val]) Len() int {
	return cache.list.Len()
}
