package main

import (
	"strings"
	"testing"
)

func BenchmarkTeeReader(b *testing.B) {
	repeat := 300
	r1 := strings.NewReader(strings.Repeat("Daje Roma ", repeat))
	l := r1.Len()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashAndSend(r1, l)
	}

}
