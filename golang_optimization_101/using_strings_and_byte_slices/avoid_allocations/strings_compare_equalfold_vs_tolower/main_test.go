package main

import (
	"strings"
	"testing"
)

var ss = []string{
	"AbcDefghijklmnOpQrStUvwxYz1234567890",
	"abcDefghijklmnopQRSTuvwXYZ1234567890",
	"aBcDefgHIjklMNOPQRSTuvwxyz1234567890",
}

func Benchmark_EqualFold(b1 *testing.B) {
	for i := 0; i < b1.N; i++ {
		for _, a := range ss {
			for _, b := range ss {
				r := strings.EqualFold(a, b)
				if !r {
					panic("not equal!")
				}
			}
		}
	}
}
func Benchmark_CompareToLower(b1 *testing.B) {
	for i := 0; i < b1.N; i++ {
		for _, a := range ss {
			for _, b := range ss {
				r := strings.ToLower(a) ==
					strings.ToLower(b)
				if !r {
					panic("not equal!")
				}
			}
		}
	}
}
