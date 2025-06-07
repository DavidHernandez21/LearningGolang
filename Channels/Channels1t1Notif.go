package main

import (
	"fmt"
	"sort"
)

func main2() {
	// values := make([]byte, 32*1024*1024)
	values := []int32{1, -4, 23, 325, -32, 5}
	// if _, err := rand.Read(values); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	fmt.Println(values)

	done := make(chan struct{}) // can be buffered or not

	// The sorting goroutine
	go func() {
		sort.Slice(values, func(i, j int) bool {
			return values[i] < values[j]
		})
		// Notify sorting is done.
		fmt.Println("Just before notification")
		done <- struct{}{}
	}()

	// do some other things ...

	<-done // waiting here for notification

	fmt.Println(values[0], values[len(values)-1])
}
