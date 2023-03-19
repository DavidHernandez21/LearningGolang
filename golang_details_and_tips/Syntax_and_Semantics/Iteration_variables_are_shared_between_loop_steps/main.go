package main

func loop1(s []int) []*int {
	r := make([]*int, len(s))
	for i, v := range s {
		r[i] = &v
	}
	return r
}
func loop2(s []int) []*int {
	r := make([]*int, len(s))
	for i := range s {
		// v := s[i]
		// r[i] = &v
		r[i] = &s[i]
	}
	return r
}
func printAll(s []*int) {
	for i := range s {
		print(*s[i])
	}
	println()
}
func main() {
	var s1 = []int{1, 2, 3}
	printAll(loop1(s1)) // 333
	var s2 = []int{1, 2, 3}
	printAll(loop2(s2)) // 123
}
