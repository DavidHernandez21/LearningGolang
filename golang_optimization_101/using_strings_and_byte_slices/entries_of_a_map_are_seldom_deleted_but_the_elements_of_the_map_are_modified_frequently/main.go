package main

import t "testing"

var key = []byte{'k', 'e', 'y'}
var m1 = map[string]int{"key": 0}

func modify1() {
	m1[string(key)]++
	// (logically) equivalent to:
	// m1[string(key)] = m1[string(key)] + 1
}

var m2 = map[string]*int{"key": new(int)}

func modify2() {
	*m2[string(key)]++
	// equivalent to:
	// p := m2[string(key)]; *p = *p + 1
}
func main() {
	stat := func(f func()) int {
		allocs := t.AllocsPerRun(1, f)
		return int(allocs)
	}
	println(stat(modify1)) // 1
	println(stat(modify2)) // 0
}
