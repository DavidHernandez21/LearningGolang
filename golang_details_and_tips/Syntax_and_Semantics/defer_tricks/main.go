package main

type T struct{}

func (t T) M(n int) T {
	print(n)
	return t
}
func main() {
	var t T
	defer t.M(1).M(5).M(2)
	t.M(3)
}
