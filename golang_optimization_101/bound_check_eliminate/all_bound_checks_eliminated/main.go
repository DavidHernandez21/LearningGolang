package main

func f2a(s []int) {
	for i := range s {
		_ = s[i]
		_ = s[i:]
		_ = s[:i+1]
	}
}
func f2b(s []int) {
	for i := 0; i < len(s); i++ {
		_ = s[i]
		_ = s[i:]
		_ = s[:i+1]
	}
}
func f2c(s []int) {
	for i := len(s) - 1; i >= 0; i-- {
		_ = s[i]
		_ = s[i:]
		_ = s[:i+1]
	}
}
func f2d(s []int) {
	for i := len(s); i > 0; {
		i--
		_ = s[i]
		_ = s[i:]
		_ = s[:i+1]
	}
}
func f2e(s []int) {
	for i := 0; i < len(s)-1; i += 2 {
		_ = s[i]
		_ = s[i:]
		_ = s[:i+1]
	}
}
func main() {}
