package main

import "runtime"

//go:noinline
func f(i int) byte {
	var a [1 << 13]byte // allocated on stack and make stack grow
	return a[i]
}

// stack frame size go run -gcflags=-S .\stacks_grows_shrinkage.go
func main() {
	var x int
	println(&x)  // <address 1>
	f(1)         // (make stack grow)
	println(&x)  // <address 2>
	runtime.GC() // (make stack shrink)
	println(&x)  // <address 3>
	runtime.GC() // (make stack shrink)
	println(&x)  // <address 4>
	runtime.GC() // (stack does not shrink)
	println(&x)  // <address 4>
}
