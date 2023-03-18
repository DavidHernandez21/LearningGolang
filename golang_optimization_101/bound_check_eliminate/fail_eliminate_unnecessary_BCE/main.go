package main

func fa(s []byte) {
	if len(s) >= 2 {
		r := len(s) % 2
		_ = s[r] // Found IsInBounds
	}
}
func fb(s, x, y []byte) {
	n := copy(s, x)
	copy(s[n:], y) // Found IsSliceInBounds
	_ = x[n:]      // Found IsSliceInBounds
}
func fc(s []byte) {
	const N = 6
	for i := 0; i < len(s)-(N-1); i += N {
		_ = s[i+N-1] // Found IsInBounds
	}
}
func fd(data []int, check func(int) bool) []int {
	var k = 0
	for _, v := range data {
		if check(v) {
			data[k] = v // Found IsInBounds
			k++
		}
	}
	return data[:k] // Found IsSliceInBounds
}
