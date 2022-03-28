package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sys/cpu"
)

type IntVars struct {
	i1 int64
	i2 int64
}

func IncrementManyTimes(val *int64, times int) *int64 {
	for i := 0; i < times; i++ {
		*val++
	}
	return val
}

type IntVarsPadded struct {
	i1 int64
	_  cpu.CacheLinePad // padding
	i2 int64
}

func IncrementParallel(val1, val2 *int64, times int) (<-chan *int64, <-chan *int64) {
	// define waitgroup, so we can wait until background incrementing tasks are finished
	// wg := &sync.WaitGroup{}
	// wg.Add(2)
	out1 := make(chan *int64, 1)
	out2 := make(chan *int64, 1)

	// run first incrementing task in background
	go func() {
		defer close(out1)
		out1 <- IncrementManyTimes(val1, times)
		// out1 <- IncrementManyTimes(val2, times)

		// wg.Done()
	}()

	// run second incrementing task in background
	go func() {
		defer close(out2)
		out2 <- IncrementManyTimes(val2, times)
		// wg.Done()
	}()

	// wait for tasks to complete
	// // wg.Wait()
	// close(out1)
	// close(out2)
	return out1, out2
}

func merge(channels ...<-chan *int64) <-chan *int64 {
	var wg sync.WaitGroup
	wg.Add(len(channels))
	out := make(chan *int64)

	for _, ch := range channels {
		go func(ch <-chan *int64) {
			defer func() {
				wg.Done()
			}()
			for val := range ch {
				out <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)

	}()
	return out
}

func StartTimer(name string) func() {
	t := time.Now()
	fmt.Println(name, "started")

	return func() {
		// d := time.Now().Sub(t)
		d := time.Since(t)
		fmt.Println(name, "took", d)
	}

}

func main() {
	timer := StartTimer("main padded")
	defer timer()
	a := IntVarsPadded{}
	// fmt.Println(a.i1)
	// fmt.Println(a.i2)
	for val := range merge(IncrementParallel(&a.i1, &a.i2, 100000000)) {
		fmt.Println(*val)
	}
	// ch1, ch2 := IncrementParallel(&a.i1, &a.i2, 100000000)
	// v1, v2 := <-ch1, <-ch2
	// fmt.Println(*v1, *v2)
	// for val := range IncrementParallel(&a.i1, &a.i2, 100000000) {
	// 	fmt.Println(*val)
	// }
}
