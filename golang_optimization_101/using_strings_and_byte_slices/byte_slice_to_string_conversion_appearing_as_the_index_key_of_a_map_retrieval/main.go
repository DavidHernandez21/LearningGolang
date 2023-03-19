package main

import (
	t "testing"
)

var m = map[string]int{}
var key = []byte{'k', 'e', 'y'}
var n int

func get() {
	n = m[string(key)]
	// fmt.Println(n)
}
func inc() {
	m[string(key)]++
}

func set() {
	m[string(key)] = 123
}

// This optimization also works if the key presents
// as a struct or array composite literal form
type T struct {
	a int
	b bool
	k [2]string
}

var m_struct = map[T]int{}
var key_struct = []byte{'k', 'e', 'y', 99: 'z'}
var n_struct int

func get_struct() {
	n = m_struct[T{k: [2]string{1: string(key_struct)}}]
}

func main() {
	stat := func(f func()) int {
		allocs := t.AllocsPerRun(1, f)
		return int(allocs)
	}

	println(stat(set))        // 1
	println(stat(inc))        // 1
	println(stat(get))        // 0
	println(stat(get_struct)) // 0
	// fmt.Println([5]int{3: 2})
	// println(m)
	// fmt.Println(m)
}
