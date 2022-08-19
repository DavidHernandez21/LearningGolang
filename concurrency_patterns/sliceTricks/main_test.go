package main

import "testing"

const (
	SliceLenght = 1000
	NumWorkers  = 3
)

func BenchmarkMySplitSlice(b *testing.B) {

	nums := make([]int, SliceLenght)
	for i := 0; i < SliceLenght; i++ {
		nums[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for num := range mySplitSlice(NumWorkers, nums...) {
			_ = num
		}
	}

}

func BenchmarkBatchesFromSlice(b *testing.B) {

	nums := make([]int, SliceLenght)
	for i := 0; i < SliceLenght; i++ {
		nums[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for num := range batchesFromSlice(NumWorkers, nums...) {
			_ = num
		}
	}

}
