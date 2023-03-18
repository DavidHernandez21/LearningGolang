package main

import (
	"testing"
)

var s = []byte{32: 'b'} // len(s) == 33
var r string

func Concat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r = string(s) + string(s)
	}
}

// There are 3 allocations made within the Concat function. Two of them are caused by the byte slice
// to string conversions string(s), and the sizes of the two memory blocks carrying the underlying
// bytes of the two result strings are both 48 (which is the smallest size class which is not smaller than 33).
// The third allocation is caused by the string concatenation, and the size of the result memory
// block is 80 (the smallest size class which is not smaller than 66). The three allocations allocate 176
// (48+48+80) bytes totally. In the end, 14 bytes are wasted. And 44 (15 + 15 + 14) bytes are wasted
// during executing the Concat function.

func main() {
	// fmt.Println(s)
	br := testing.Benchmark(Concat)
	println(br.AllocsPerOp())       // 3
	println(br.AllocedBytesPerOp()) // 176
}
