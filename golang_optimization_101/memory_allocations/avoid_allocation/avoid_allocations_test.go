package main

import "testing"

func buildOrginalData() []int {
	s := make([]int, 10)
	for i := range s {
		s[i] = i
	}
	return s
}
func check(v int) bool {
	return v%2 == 0
}
func FilterOneAllocation(data []int) []int {
	var r = make([]int, 0, len(data))
	for _, v := range data {
		if check(v) {
			r = append(r, v)
		}
	}
	return r
}

func FilterNoAllocations(data []int) []int {

	var k = 0

	for i, v := range data {
		if check(v) {

			// Is this really needed?
			data[i] = data[k]

			data[k] = v

			k++
		}
	}
	// fmt.Println(data)
	return data[:k]
}

// func main() {
// 	data := buildOrginalData()
// 	// fmt.Printf("%p\n", data)
// 	// fmt.Println(FilterOneAllocation(data))
// 	data = FilterNoAllocations(data)
// 	fmt.Printf("%v\n", data)

// }

func Benchmark_FilterOneAllocation(b *testing.B) {
	data := buildOrginalData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FilterOneAllocation(data)
	}
}

func Benchmark_FilterNoAllocations(b *testing.B) {
	data := buildOrginalData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FilterNoAllocations(data)
	}
}
