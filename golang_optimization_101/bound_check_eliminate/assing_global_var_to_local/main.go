package main

var s = make([]int, 5)

func fa0() {
	for i := range s {
		s[i] = i // Found IsInBounds
	}
}
func fa1() {
	s := s
	for i := range s {
		s[i] = i
	}
}
func fa2(x []int) {
	for i := range x {
		x[i] = i
	}
}

func main() {}
