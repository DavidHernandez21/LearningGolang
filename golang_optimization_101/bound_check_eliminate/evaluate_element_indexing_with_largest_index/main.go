package main

import "fmt"

func f3a(s []int32) int32 {
	return s[0] | // Found IsInBounds (line 5)
		s[1] | // Found IsInBounds
		s[2] | // Found IsInBounds
		s[3] // Found IsInBounds
}
func f3b(s []int32) int32 {
	return s[3] | // Found IsInBounds (line 12)
		s[0] |
		s[1] |
		s[2]
}
func main() {
	s := []int32{1, 2, 3, 4}
	fmt.Println(f3b(s))
}
