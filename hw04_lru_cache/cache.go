package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type keyValuePair struct {
	key   Key
	value interface{}
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem

	mtx sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mtx.Lock()
	defer lc.mtx.Unlock()

	item, found := lc.items[key]
	if found {
		lc.queue.PushFront(item)
		item.Value = keyValuePair{key, value}
		return true
	}

	for lc.queue.Len() >= lc.capacity {
		itemToDelete := lc.queue.Back()
		keyToDelete := itemToDelete.Value.(keyValuePair).key
		delete(lc.items, keyToDelete)
		lc.queue.Remove(itemToDelete)
	}
	item = lc.queue.PushFront(keyValuePair{key, value})
	lc.items[key] = item
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mtx.RLock()
	defer lc.mtx.RUnlock()

	item, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(item)
		return item.Value.(keyValuePair).value, true
	}

	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mtx.Lock()
	defer lc.mtx.Unlock()

	clear(lc.items)

	lc.queue = NewList()
}
