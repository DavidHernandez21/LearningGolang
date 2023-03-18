package pointers

import "testing"

func f(sum *int, s []int) {
	for _, v := range s { // line 7
		*sum += v // line 8
	}
}
func g(sum *int, s []int) {
	var n = 0
	for _, v := range s { // line 14
		n += v // line 15
	}
	*sum = n
}

var s = make([]int, 1024)
var r int

// The outputted assembly instructions show the pointer sum
// is dereferenced within the loop in the
// f function. A dereference operation is a memory operation.
// For the g function, the dereference
// operations happen out of the loop, and the instructions
// generated for the loop only process registers.
// It is much faster to let CPU instructions process
// registers than process memory, which is why the g
// function is much more performant than the f function.

func Benchmark_f(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f(&r, s)
	}
}
func Benchmark_g(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g(&r, s)
	}
}
