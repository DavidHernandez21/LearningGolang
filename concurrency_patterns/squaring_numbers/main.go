package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for k := range nums {
			out <- nums[k]
		}
		close(out)
	}()
	return out
}

// When the number of values to be sent is known at channel creation time
//
//	a buffer can simplify the code.
func gen2(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	for k := range nums {
		out <- nums[k]
	}
	close(out)
	return out
}

func sq(done chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				fmt.Println("done sq")
				return
			}
		}
	}()
	return out
}

func merge(done chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				fmt.Println("done merge")
				return
			}
		}
	}

	wg.Add(len(cs))
	for k := range cs {
		go output(cs[k])
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	SliceLenght := 200
	nums := make([]int, SliceLenght)
	for i := 0; i < SliceLenght; i++ {
		nums[i] = i
	}

	in := gen2(nums...)

	// Distribute the sq work across numWorkers goroutines that both read from in.
	numWorkers := 5
	workers := make([]<-chan int, numWorkers)

	done := make(chan struct{})
	defer close(done)

	for i := 0; i < numWorkers; i++ {
		workers[i] = sq(done, in)
	}
	// c1 := sq(in)
	// c2 := sq(in)

	// Consume the merged output from workers.
	for n := range merge(done, workers...) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}
