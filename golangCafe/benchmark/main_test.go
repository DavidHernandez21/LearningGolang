package main

import (
	"fmt"
	"testing"
)

func TestFibonacci(t *testing.T) {
	testCases := []struct{ seqNo, expected int }{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{5, 5},
		{10, 55},
	}
	for n, tc := range testCases {
		testCase := tc
		var m = make(map[int]int)
		t.Run(fmt.Sprintf("test-case-%d", n), func(t *testing.T) {
			t.Parallel()
			out := quickFibonacciNoVarDeclaration(testCase.seqNo, m)
			if out != testCase.expected {
				t.Errorf("fibonacci(%d) = %d; got %d", testCase.seqNo, testCase.expected, out)
			}
		})
	}
}

// func BenchmarkFibonacci(b *testing.B) {
// 	seqs := []int{10, 20, 30, 40}
// 	for _, seq := range seqs {
// 		b.Run(fmt.Sprintf("benchmark-seq-%d", seq), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				fibonacci(seq)
// 			}
// 		})
// 	}
// }

func BenchmarkQuickFibonacci(b *testing.B) {
	seqs := []int{10, 20, 30, 40}
	for _, seq := range seqs {
		var m = make(map[int]int)
		b.Run(fmt.Sprintf("benchmark-seq-%d", seq), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				quickFibonacci(seq, m)
			}
		})
	}
}

func BenchmarkQuickFibonacciNoVarDeclaration(b *testing.B) {
	seqs := []int{10, 20, 30, 40}
	for _, seq := range seqs {
		var m = make(map[int]int)
		b.Run(fmt.Sprintf("benchmark-seq-%d", seq), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				quickFibonacciNoVarDeclaration(seq, m)
			}
		})
	}
}
