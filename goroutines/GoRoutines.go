// Description: This program demonstrates the use of goroutines and WaitGroups in Go.
package main

import (
	"fmt"
	"sync"
)

func sayHello(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello from goroutine")
}

func main1() {
	var wg = sync.WaitGroup{}

	arr := [4]uint16{1, 5, 7, 9}

	wg.Add(len(arr))

	for i := range arr {
		go sayHello(&wg)
		fmt.Printf("cycle number #%d\t number: #%d\n", i, arr[i])
	}

	wg.Wait()
	fmt.Println("main goroutine exit")
}
