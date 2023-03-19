package main

import "testing"

var v = 9999999
var x, y interface{}

func Benchmark_BoxBox(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x = v // needs one allocation
		y = v // needs one allocation
	}
}
func Benchmark_BoxAssign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x = v // needs one allocation
		y = x // no allocations
	}
}
