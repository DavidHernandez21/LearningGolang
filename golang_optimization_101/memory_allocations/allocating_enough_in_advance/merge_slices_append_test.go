package allocations

import "testing"

func getData() [][]int {
	return [][]int{
		{1, 2},
		{9, 10, 11},
		{6, 2, 3, 7},
		{11, 5, 7, 12, 16},
		{8, 5, 6},
	}
}

func MergeWithOneLoop(data ...[]int) []int {
	var r []int
	for _, s := range data {
		r = append(r, s...)
	}
	return r
}
func MergeWithTwoLoops(data ...[]int) []int {
	n := 0
	// handle the integer overflowing case
	for _, s := range data {
		if k := n + len(s); k < n {
			panic("slice length overflows")
		} else {
			n = k
		}
	}

	r := make([]int, 0, n)
	for _, s := range data {
		r = append(r, s...)
	}
	return r
}
func Benchmark_MergeWithOneLoop(b *testing.B) {
	data := getData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MergeWithOneLoop(data...)
	}
}
func Benchmark_MergeWithTwoLoops(b *testing.B) {
	data := getData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MergeWithTwoLoops(data...)
	}
}
