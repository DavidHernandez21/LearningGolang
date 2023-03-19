package main

import "testing"

var m = map[int]int{}

func BenchmarkIncrement(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m[99]++
	}
}
func BenchmarkPlusOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m[99] += 1
	}
}
func BenchmarkAddition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m[99] = m[99] + 1
	}
}
