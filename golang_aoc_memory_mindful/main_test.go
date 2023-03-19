package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func Example() {
	file, err := os.Open("testdata/input.txt")
	if err != nil {
		log.Fatalf("could not open input file: %v", err)
	}
	defer file.Close()

	if err := Solve(file, os.Stdout); err != nil {
		log.Fatalf("could not solve: %v", err)
	}
	// Output: 2518
}

func Benchmark(b *testing.B) {
	input, err := os.ReadFile("testdata/input.txt")
	if err != nil {
		b.Fatalf("could not read input file: %v", err)
	}

	r := bytes.NewReader(input)
	w := io.Discard

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.Reset(input)
		_ = Solve(r, w)
	}
}
