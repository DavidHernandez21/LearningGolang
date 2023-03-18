package main

import (
	"testing"
)

var s = []byte{0: '$', 32: 'x'} // len(s) == 33
var s32 = []byte{0: '$', 31: 'x'}

// expression doesn’t allocate if at least one of concatenated
// operands is a non-blank string constant
func f() string {
	return (" " + string(s) + string(s))[1:]
}

func f32() string {
	return (" " + string(s32) + string(s32))[1:]
}
func g() string {
	return string(s) + string(s)
}

func g32() string {
	return string(s32) + string(s32)
}

var x string

func stat(add func() string) int {
	c := func() {
		x = add()
	}
	allocs := testing.AllocsPerRun(1, c)

	return int(allocs)
}

func main() {
	println(stat(f)) // 1
	println(stat(g)) // 3
	println(stat(f32))
	println(stat(g32))
	// fmt.Println(s)
	// fmt.Println(len(f()), len(g()))
}
