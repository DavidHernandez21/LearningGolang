package main

import "unsafe"

type T1 struct {
	a int8
	b int64
	c int16
}
type T2 struct {
	a int8
	c int16
	b int64
}

// • The alignment guarantee of the struct type is the same as its largest alignment guarantee of
// its filed types. Here is the alignment guarantee (8, a native word) of type int64. This means
// the distance between the address of the field b and a of a value of the struct type is a multiple
// of 8. Clever compilers should choose the minimum possible value: 8. To get the desired
// alignment, 7 bytes are padded after the field a.
// • The size of the struct type must be a multiple of the alignment guarantee of the struct type. So
// considering the existence of the field c, the minimum possible size is 24 (8x3), which should
// be used by clever compilers. To get the desired size, 6 bytes are padded after the field c

// In practice, generally, we should make related fields adjacent to get good readability, and only order
// fields in the most memory saving way when it really needs.

func main() {
	// The printed values are got on
	// 64-bit architectures.
	println(unsafe.Sizeof(T1{})) // 24
	println(unsafe.Sizeof(T2{})) // 16
}
