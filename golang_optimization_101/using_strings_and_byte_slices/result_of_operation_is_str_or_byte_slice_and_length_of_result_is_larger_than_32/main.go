package main

import "testing"

var str = "1234567890abcdef" // len(str) == 16
func f() {
	x := str + str          // does not escape
	y := []byte(x)          // does not escape
	println(len(y), cap(y)) // 32 32
	z := string(y)          // does not escape
	println(len(x), len(z)) // 32 32
}
func g() {
	x := str + str + "x"
	y := []byte(x)
	println(len(y), cap(y)) // 33 48
	z := string(y)
	println(len(x), len(z))
}
func stat(fn func()) int {
	allocs := testing.AllocsPerRun(1, fn)
	return int(allocs)
}
func main() {
	println(stat(f)) // 0
	println(stat(g)) // 3
}
