package main

import (
	"testing"
)

const (
	s  = "1234567890abcdef"
	s1 = "abcdef1234567890"
)

func f1(x, y string) {
	switch {
	case x == y: // ... handle 1
	case x < y: // ... handle 2
	default: // ... handle 3
	}
}
func f2(x, y string) {
	switch {
	case x < y: // ... handle 2
	case x == y: // ... handle 1
	default: // ... handle 3
	}
}
func f3(x, y string) {
	switch {
	case x < y: // ... handle 2
	case x > y: // ... handle 3
	default: // ... handle 1
	}
}

func BenchmarkStringsCompareEqualNotDefault(b *testing.B) {
	// s2 := strings.Repeat(s, 100)
	// s3 := strings.Repeat(s1, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f1(s, s1)
	}

}

func BenchmarkStringsCompareEqualInDefault(b *testing.B) {
	// s2 := strings.Repeat(s, 100)
	// s3 := strings.Repeat(s1, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f3(s, s1)
	}

}
