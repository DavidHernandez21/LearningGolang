package main

import "fmt"

func mySplitSlice(num_workers int, nums ...int) [][]int {
	l := len(nums)
	chunk_size := l / num_workers
	// fmt.Println(chunk_size)

	// var chunks [][]int
	// wg.Add()
	var end int
	chunkOfChunks := make([][]int, num_workers+l%num_workers)

	for k := range chunkOfChunks {
		chunkOfChunks[k] = make([]int, 0, chunk_size)
	}

	n := 0
	for i := 0; i < l; i += chunk_size {

		end = i + chunk_size

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > l {
			end = l
		}

		// copy(chunks, nums[i:end])
		// copy(chunkOfChunks[n], nums[i:end])
		chunkOfChunks[n] = nums[i:end]
		n++

	}
	// fmt.Println(n, len(chunkOfChunks), l%num_workers)
	return chunkOfChunks
}

func batchesFromSlice(batchSize int, nums ...int) [][]int {

	batches := make([][]int, 0, (len(nums)+batchSize-1)/batchSize)

	for batchSize < len(nums) {
		nums, batches = nums[batchSize:], append(batches, nums[0:batchSize:batchSize])
	}
	batches = append(batches, nums)

	return batches

}

func main() {

	SliceLenght := 20
	nums := make([]int, SliceLenght)
	for i := 0; i < SliceLenght; i++ {
		nums[i] = i
	}
	nums2 := make([]int, SliceLenght)
	copy(nums2, nums)
	// actions := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	batchSize := 3

	fmt.Println(batchesFromSlice(batchSize, nums...))

	fmt.Println(mySplitSlice(batchSize, nums2...))
	// fmt.Println(nums2)

}
