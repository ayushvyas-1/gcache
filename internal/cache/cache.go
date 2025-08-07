package cache

import (
	"fmt"
	"sync"
)

type LRUCache struct {
	capacity int
	cache    map[string]*DoublyNode
	list     *DoublyLinkedList
	mu       sync.RWMutex
}

type CacheItem struct {
	key   string
	value string
}

func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		panic("LRUCache capacity must be greater than 0")
	}

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*DoublyNode),
		list:     NewDoublyLinkedList(),
	}
}

func (lru *LRUCache) Get(key string) (string, bool) {

	lru.mu.Lock()         // mutex lock -- blocks RW
	defer lru.mu.Unlock() // unlocks when the func end

	if node, exist := lru.cache[key]; exist {

		lru.list.Remove(node)
		lru.list.InsertAtFront(node)

		item := node.GetData().(*CacheItem)
		return item.value, true
	}

	return "", false
}

func (lru *LRUCache) Put(key, value string) {
	lru.mu.Lock()         // mutex lock -- blocks RW
	defer lru.mu.Unlock() // unlocks when the func end

	if node, exists := lru.cache[key]; exists {

		item := node.GetData().(*CacheItem)
		item.value = value

		lru.list.Remove(node)
		lru.list.InsertAtFront(node)
		return
	}

	if lru.list.Count() == lru.capacity {
		tail := lru.list.last
		if tail != nil {
			lru.list.Remove(tail)
			evicted := tail.GetData().(*CacheItem)
			delete(lru.cache, evicted.key)
		}
	}

	item := &CacheItem{key, value}

	node := NewDoublyNode(item)
	lru.list.InsertAtFront(node)
	lru.cache[key] = node
}

func (lru *LRUCache) Delete(key string) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if node, exists := lru.cache[key]; exists {
		lru.list.Remove(node)
		delete(lru.cache, key)
		return true
	}
	return false
}

func (lru *LRUCache) Size() int {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	return lru.list.Count()
}

func (lru *LRUCache) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.cache = make(map[string]*DoublyNode)
	lru.list = NewDoublyLinkedList()
}

func (lru *LRUCache) Contains(key string) bool {
	lru.mu.RLock()
	defer lru.mu.RUnlock()

	_, exists := lru.cache[key]
	return exists
}

//for quick look

func (lru *LRUCache) Print() {

	lru.mu.RLock()
	defer lru.mu.RUnlock()

	fmt.Printf("Cache state (MRU->LRU, size: %d/%d):\n", lru.list.Count(), lru.capacity)

	current := lru.list.head
	for current != nil {
		item := current.GetData().(*CacheItem)
		fmt.Printf("[%s : %s] ", item.key, item.value)
		current = current.next
	}

	fmt.Println()

}
