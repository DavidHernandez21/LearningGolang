package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	capacity := 10
	s := make([]int, 0, capacity)

	wg.Add(capacity)

	for i := 0; i < capacity; i++ {
		go func(i int) {
			defer wg.Done()
			s = append(s, i)
		}(i)
	}

	wg.Wait()
	fmt.Println(s)

}
