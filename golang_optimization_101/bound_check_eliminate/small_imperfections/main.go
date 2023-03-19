package main

func f2c(s []int) {
	for i := len(s) - 1; i >= 0; i-- {
		//_ = s[i]
		_ = s[:i+1] // line 7: Found IsSliceInBounds
	}
}
func f2e(s []int) {
	for i := 0; i < len(s)-1; i += 2 {
		//_ = s[i]
		_ = s[:i+1] // line 14: Found IsSliceInBounds
	}
}

// We may give the compiler some hints by turning on the comment lines to remove these bound
// checks.

func main() {}
