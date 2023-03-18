package main

import "strconv"

func f(x int) string { // x escapes to heap
	if x >= 0 && x < 10 {
		return "0123456789"[x : x+1]
	}
	return g(&x)
}
func g(x *int) string {
	escape(x) // for demo purpose
	return strconv.Itoa(*x)
}

func f_optimized(x int) string { // x does not escape to heap
	if x >= 0 && x < 10 {
		return "0123456789"[x : x+1]
	}

	x2 := x
	return g(&x2)
}

var sink interface{}

//go:noinline
func escape(x interface{}) {
	sink = x
	sink = nil
}
func main() {
	var a = f(100)
	println(a)
	var b = f_optimized(100)
	println(b)
}
