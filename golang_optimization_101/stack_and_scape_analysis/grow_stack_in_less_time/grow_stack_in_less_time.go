package main

import (
	"fmt"
	"time"
)

func demo(n int) byte {
	var a [8192]byte
	var b byte
	if n--; n > 0 {
		b = demo(n)
	}
	return a[n] + b
}
func foo(c chan time.Duration) {
	start := time.Now()
	demo(8192)
	c <- time.Since(start)
}
func bar(c chan time.Duration) {
	start := time.Now()
	// Use a dummy anonymous function to enlarge stack.
	func(x *interface{}) {
		type _ int // avoid being inlined
		if x != nil {
			*x = [1024 * 1024 * 64]byte{}
		}
	}(nil)
	demo(8192)
	c <- time.Since(start)
}

// The reason is simple, only one stack growth happens in the lifetime of the bar goroutine,
// whereas more than 10 stack growths happen in the lifetime of the foo goroutine.
func main() {
	var c = make(chan time.Duration, 1)
	go foo(c)
	fmt.Println("foo:", <-c)
	go bar(c)
	fmt.Println("bar:", <-c)
}
