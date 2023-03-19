package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func Solve(r io.Reader, w io.Writer) error {
	datastream, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	position, err := blockPosition(datastream, 14)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%d", position)
	return err
}

func blockPosition(datastream []byte, size int) (int, error) {
	for i := 0; i < len(datastream)-size; i++ {
		if !hasDuplicates(datastream[i : i+size]) {
			return i + size, nil
		}
	}

	return 0, errors.New("not found")
}

func hasDuplicates(block []byte) bool {
	// seen := make(map[byte]bool)
	// seen := [256]bool{}
	seen := make([]uint8, 256)
	for i := range block {
		if seen[block[i]] == 1 {
			return true
		}
		seen[block[i]] = 1
	}

	return false
}

func hasDuplicates2(block []byte) bool {
	// seen := make(map[byte]bool)
	// seen := [256]bool{}
	seen := make([]uint8, 256)
	for i := range block {

		seen[block[i]] = block[i]
	}

	sort.Slice(seen, func(i, j int) bool {
		return seen[i] > seen[j]
	})

	return seen[13] == 0

	// return false
}

func main() {
	file, err := os.Open("testdata/input.txt")
	if err != nil {
		log.Fatalf("could not open input file: %v", err)
	}
	defer file.Close()

	if err := Solve(file, os.Stdout); err != nil {
		log.Fatalf("could not solve: %v", err)
	}

}
