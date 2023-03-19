package main

import "fmt"

func getData() [][]int {
	return [][]int{
		{1, 2},
		{9, 10, 11},
		{6, 2, 3, 7},
		{11, 5, 7, 12, 16},
		{8, 5, 6},
	}
}

const format = "Allocate from %v to %v (when append slice#%v).\n"

const cap_format = "oldCap: %v, newCap: %v, slice#%v.\n"

func MergeWithOneLoop(data [][]int) []int {
	var oldCap int
	var r []int
	for i, s := range data {
		r = append(r, s...)
		fmt.Printf(cap_format, oldCap, cap(r), i)
		if oldCap == cap(r) {
			continue
		}
		fmt.Printf(format, oldCap, cap(r), i)
		oldCap = cap(r)
	}
	return r
}
func main() {
	MergeWithOneLoop(getData())
}
