package main

import "log"

// If a slice is short-lived, then we could allocate it
//  with an estimated large enough capacity. There
// might be some memory wasted temporarily,
// but the element memory will be released soon. Even if
// the estimated capacity is proved to be not large enough,
//  there will still be several allocations saved.

func main() {
	x := make([]int, 100, 500)
	// println(cap(x[:len(x):len(x)]))
	for i := 0; i < len(x); i++ {
		x[i] = i
	}
	y := make([]int, 500)
	for i := 100; i < len(y); i++ {
		y[i] = i
	}
	// log.Println(y)
	a := append(x, y...)
	b := append(x[:len(x):len(x)], y...)
	println(cap(a)) // 1024
	println(cap(b)) // 640
	log.Println(len(a))
	log.Println(len(b))
}
