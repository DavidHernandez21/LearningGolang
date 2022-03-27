package main

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
)

func main() {
	cache, err := lru.New(128)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 256; i++ {
		cache.Add(i, i)
	}

	fmt.Println(cache.Get(132))
	fmt.Println(cache.Contains(132))
	fmt.Println(cache.Peek(132))
	k, v, ok := cache.GetOldest()
	fmt.Println(v, ok)
	fmt.Println(cache.GetOldest())
	cache.Remove(k)
	fmt.Println(cache.RemoveOldest())
	fmt.Println(cache.Keys())
	fmt.Println(cache.Len())
	cache.Purge()
	fmt.Println(cache.Len())
}
