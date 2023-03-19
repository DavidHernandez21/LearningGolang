package main

import (
	"bytes"
	t "testing"
)

func verbose(x, y, z []byte) {
	switch {
	case string(x) == string(y):
	// do something
	case string(x) == string(z):
		// do something
	}
}
func clean(x, y, z []byte) {
	switch string(x) {
	case string(y):
	// do something
	case string(z):
		// do something
	}
}

// Note, two branches are enough
// to form a three-way comparison.
// The bytes.Compare function way is often more performant
// for the cases in which THREE-WAY comparisons are needed
func doSomething(x, y []byte) {
	switch bytes.Compare(x, y) {
	case -1:
	// ... do something 1
	case 1:
	// ... do something 2
	default:
		// ... do something 3
	}
}

func main() {
	x := []byte{1023: 'x'}
	y := []byte{1023: 'y'}
	z := []byte{1023: 'z'}
	stat := func(f func(x, y, z []byte)) int {
		allocs := t.AllocsPerRun(1, func() {
			f(x, y, z)
		})
		return int(allocs)
	}
	// fmt.Println(len(x), len(string(x)))
	println(stat(verbose)) // 0
	println(stat(clean))   // 3
}
