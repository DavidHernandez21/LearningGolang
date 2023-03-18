package main

func f1a(s []struct{}, index int) {
	_ = s[index] // line 5: Found IsInBounds
	_ = s[index]
	_ = s[index:]
	_ = s[:index+1]
}
func f1b(s []byte, index int) {
	s[index-1] = 'a' // line 12: Found IsInBounds
	_ = s[:index]
}

func f1c(a [5]int) {
	_ = a[0]
	_ = a[4]
}
func f1d(s []int) {
	if len(s) > 2 {
		_, _, _ = s[0], s[1], s[2]
	}
}
func main() {}
