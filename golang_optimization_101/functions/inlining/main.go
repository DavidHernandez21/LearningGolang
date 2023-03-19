package main

import "fmt"

func bar(a, b int) int {
	return a*a - b*b + 2*(a-b)
}
func foo(x, y int) int {
	var a = bar(x, y)
	var b = bar(y, x)
	var c = bar(a, b)
	var d = bar(b, a)
	return c*c + d*d
}

func main() {
	fmt.Println(foo(1, 2))
}
