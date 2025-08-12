package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ayushvyas-1/gcache/internal/cache"
)

func TestBasicOperations(t *testing.T) {
	cache := cache.NewLRUCache(3)

	cache.Put("a", "1")
	cache.Put("b", "2")
	cache.Put("c", "3")

	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	if val, ok := cache.Get("a"); !ok || val != "1" {
		t.Errorf("Expected 'a':'1', got '%s':%t", val, ok)
	}

	if val, ok := cache.Get("b"); !ok || val != "2" {
		t.Errorf("Expected 'b':'2', got '%s':%t", val, ok)
	}

	if val, ok := cache.Get("x"); ok {
		t.Errorf("Expected key 'x' to not exist, but got '%s'", val)
	}

}

func TestLRUEviction(t *testing.T) {
	cache := cache.NewLRUCache(2)

	cache.Put("a", "1")
	cache.Put("b", "2")
	cache.Put("c", "3") // 'a' should get evicted

	if _, ok := cache.Get("a"); ok {
		t.Error("Expected 'a' to be evicted")
	}

	if _, ok := cache.Get("b"); !ok {
		t.Error("Expected 'b' to exist")
	}

	if _, ok := cache.Get("c"); !ok {
		t.Error("Expected 'c' to exist")
	}

}

func TestClear(t *testing.T) {
	cache := cache.NewLRUCache(3)

	cache.Put("a", "1")
	cache.Put("b", "2")
	cache.Put("c", "3")

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", cache.Size())
	}

	if _, ok := cache.Get("a"); ok {
		t.Error("Expected cache to be empty after clear")
	}
}

// test for mutex locks
func TestConcurrency(t *testing.T) {
	cache := cache.NewLRUCache(100)
	var wg sync.WaitGroup

	numWorkers := 10

	numOps := 100

	wg.Add(numWorkers)

	//test concurrent writes
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := fmt.Sprintf("key_%d_%d", workerID, j)
				value := fmt.Sprintf("value_%d_%d", workerID, j)
				cache.Put(key, value)
			}
		}(i)
	}
	wg.Wait()

	wg.Add(numWorkers * 2)
	//test concurrent reads and writes

	//writers
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := fmt.Sprintf("read_key_%d_%d", workerID, j)
				value := fmt.Sprintf("read_value_%d_%d", workerID, j)
				cache.Put(key, value)
			}
		}(i)
	}

	//Readers

	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := fmt.Sprintf("read_key_%d_%d", workerID, j)
				cache.Get(key) // Don't care about result, just testing for races
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	wg.Wait()

	t.Log("Concurrency test completed successfully")

}

func TestEdgeCases(t *testing.T) {
	// Test zero capacity
	/* Self-Note :=
			here Previously test failed with segmentation fault
			because in Put(key,value) method :-
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

		if lru.list.Count() == lru.capacity {						<-- True because we never inserted any data into cache (Count = 0)
			tail := lru.list.last										and capacity = 0
			lru.list.Remove(tail)

			evicted := tail.GetData().(*CacheItem) 					<---here it called GetData() method

			GetData() => func (n *DoublyNode) GetData() interface{} {
					return n.data 									<=== there is no node!so no data!! thats why.
				}

			delete(lru.cache, evicted.key)
		}

		item := &CacheItem{key, value}

		node := NewDoublyNode(item)
		lru.list.InsertAtFront(node)
		lru.cache[key] = node
	}


	Solution was to restrict capacity >=0
	*/

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected NewLRUCache(0) to panic")
		}
	}()
	cache.NewLRUCache(0) // This should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected NewLRUCache(-1) to panic")
		}
	}()
	cache.NewLRUCache(-1)
}

func TestCapacityOne(t *testing.T) {
	// Test capacity 1 separately
	cache := cache.NewLRUCache(1)
	cache.Put("a", "1")

	if cache.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Size())
	}

	cache.Put("b", "2") // Should evict 'a'

	if _, ok := cache.Get("a"); ok {
		t.Error("Expected 'a' to be evicted in capacity-1 cache")
	}
	if val, ok := cache.Get("b"); !ok || val != "2" {
		t.Error("Expected 'b' to exist in capacity-1 cache")
	}
	if cache.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Size())
	}
}

//Benchmark tests

func BenchmarkPut(b *testing.B) {
	cache := cache.NewLRUCache(1000)

	for i := 0; b.Loop(); i++ {
		key := fmt.Sprintf("key_%d", i%500) // Some overlap to test updates
		value := fmt.Sprintf("value_%d", i)
		cache.Put(key, value)
	}
}

func BenchmarkGet(b *testing.B) {
	cache := cache.NewLRUCache(1000)

	// Pre-populate cache
	for i := 0; i < 500; i++ {
		cache.Put(fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("key_%d", i%500))
	}
}

func BenchmarkConcurrent(b *testing.B) {
	cache := cache.NewLRUCache(1000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				cache.Put(fmt.Sprintf("key_%d", i%500), fmt.Sprintf("value_%d", i))
			} else {
				cache.Get(fmt.Sprintf("key_%d", i%500))
			}
			i++
		}
	})
}
