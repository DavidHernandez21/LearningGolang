package main

import "fmt"

// import "testing"
const N = 10

func buildArray() [N]int {
	var a [N]int
	for i := 0; i < N; i++ {
		// fmt.Printf("%v, %T\n", (N-i)&i, (N-i)&i)
		a[i] = (N - i) & i
	}
	return a
}

func main() {
	// fmt.Println(buildArray())
	// a := buildArray()
	var indexTable = [10]bool{
		1: true, 2: true, 6: true, 7: true, 8: true,
	}
	fmt.Println(3 % 10)
	fmt.Println(indexTable[3%10])
	// fmt.Printf("%v, %T\n", a, a)
}
