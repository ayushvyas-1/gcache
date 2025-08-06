package main

import "fmt"

type LRUCache struct {
	capacity int
	cache    map[string]*DoublyNode
	list     *DoublyLinkedList
}

type CacheItem struct {
	key   string
	value string
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*DoublyNode),
		list:     NewDoublyLinkedList(),
	}
}

func (lru *LRUCache) Get(key string) (string, bool) {
	if node, exist := lru.cache[key]; exist {

		lru.list.Remove(node)
		lru.list.InsertAtFront(node)

		item := node.GetData().(*CacheItem)
		return item.value, true
	}

	return "", false
}

func (lru *LRUCache) Put(key, value string) {
	if node, exists := lru.cache[key]; exists {

		item := node.GetData().(*CacheItem)
		item.value = value

		lru.list.Remove(node)
		lru.list.InsertAtFront(node)
		return
	}

	if lru.list.Count() == lru.capacity {
		tail := lru.list.last
		lru.list.Remove(tail)

		evicted := tail.GetData().(*CacheItem)
		delete(lru.cache, evicted.key)
	}

	item := &CacheItem{key, value}

	node := NewDoublyNode(item)
	lru.list.InsertAtFront(node)
	lru.cache[key] = node
}

//for quick look

func (lru *LRUCache) Print() {
	fmt.Println("Cache state (MRU->LRU):")

	current := lru.list.head
	for current != nil {
		item := current.GetData().(*CacheItem)
		fmt.Printf("[%s : %s] ", item.key, item.value)
		current = current.next
	}

	fmt.Println()

}
