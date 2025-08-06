package main

import "fmt"

func main() {

	lru := NewLRUCache(9)

	lru.Put("a", "1")
	lru.Put("b", "2")
	lru.Put("c", "3")
	lru.Put("d", "4")
	lru.Put("e", "5")
	lru.Put("f", "6")
	lru.Put("g", "7")
	lru.Put("h", "8")
	lru.Put("i", "9")

	lru.Print()

	if val, ok := lru.Get("a"); ok {
		fmt.Println("GET a: ", val)
	}
	lru.Print()

	if val, ok := lru.Get("h"); ok {
		fmt.Println("GET h: ", val)
	}
	lru.Print()

	if val, ok := lru.Get("h"); ok {
		fmt.Println("GET h: ", val)
	}

	lru.Print()

	fmt.Println("add j:10")
	lru.Put("j", "10")

	lru.Print()
}
