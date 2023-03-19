package main

var s = make([]int, 256)
var a = [256]int{}

func fb1() int {
	return s[100] // Found IsInBounds
}
func fb2() int {
	return a[100]
}
func fc1(n byte) int {
	return s[n] // Found IsInBounds
}
func fc2(n byte) int {
	return a[n]
}

func main() {}
