package main

import (
	"fmt"
)

func main1() {
	// f := "ù"
	// fmt.Printf("length of f: %v", len(f))
	src := "abc"
	dst := make([]byte, 3)

	numberOfElementsCopied := copy(dst, src)
	fmt.Printf("Number Of Elements Copied: %d\n", numberOfElementsCopied)
	fmt.Printf("dst: %v\n", dst)
	fmt.Printf("src: %v\n", src)

	fmt.Printf("from []bytes to string: %s", string(dst))
}
