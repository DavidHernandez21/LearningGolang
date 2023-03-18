package main

func f9a(n int) []int {
	buf := make([]int, n+1)
	k := 0
	for i := 0; i <= n; i++ {
		buf[i] = k // Found IsInBounds
		k++
	}
	return buf
}
func f9b(n int) []int {
	buf := make([]int, n+1)
	k := 0
	for i := 0; i < len(buf); i++ {
		buf[i] = k
		k++
	}
	return buf
}

func main() {}
