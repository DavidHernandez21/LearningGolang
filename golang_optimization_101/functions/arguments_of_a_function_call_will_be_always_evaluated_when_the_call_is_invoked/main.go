package main

import (
	"log"
	"testing"
)

var debugOn = false

type stringOrInt interface {
	string | int
}

func debugPrint[T stringOrInt](s T) {
	if debugOn {
		log.Println(s)
	}
}
func main() {
	stat := func(f func()) int {
		allocs := testing.AllocsPerRun(1, f)
		return int(allocs)
	}
	var h, w string = "hello ", "world!"
	// var h, w int = 1, 2 // 1
	var n = stat(func() {
		debugPrint(h + w)
	})
	println(n) // 1
}
