package main

import (
	"fmt"
	"src/github.com/DavidHernandez21/trees/tree"
)

func main() {
	// channel := make(chan int)

	// go func(c chan int) {
	// 	v := tree.New(1).Value + 1
	// 	for i := 1; i < v; i++ {
	// 		c <- i
	// 	}
	// 	close(c)
	// }(channel)

	// var arr1 []int
	// var arr2 []int

	// fmt.Printf("%p\n", &arr1)

	// fmt.Printf("%p\n", &arr1)

	// fmt.Println([3]int{1, 2, 3} == [3]int{3, 4, 5})

	// fmt.Println(tree.New(1).Value)

	ch := make(chan int)
	done := make(chan struct{})

	t1 := tree.New(1)

	go Walk(t1, ch, done)
	go Recieve(ch, done)

	<-done
	<-done
	fmt.Println("daje")

	// go Walk(t1, ch)

	// for i := range ch {
	// 	fmt.Println(i)
	// }

	// fmt.Printf("Are the two slices equal: %t\n", Same(t1, t2))

	// fmt.Println(Same(t1, t2))

}

func Walk(t *tree.Tree, ch chan int, done chan<- struct{}) {
	v := t.Value + 1
	for i := 1; i < v; i++ {
		ch <- i
	}
	// fmt.Println("daje")
	done <- struct{}{}
	close(ch)
}

func Recieve(ch chan int, done chan<- struct{}) {
	for i := range ch {
		fmt.Println(i)
	}
	done <- struct{}{}
}

// func Same(t1, t2 *tree.Tree) bool {
// 	var slice1 []int
// 	var slice2 []int

// 	ch1 := make(chan int)
// 	ch2 := make(chan int)

// 	// fmt.Println(t1.Value, t2.Value)

// 	go Walk(t1, ch1)
// 	go Walk(t2, ch2)

// 	slice1 = AppendToSlice(slice1, ch1)
// 	slice2 = AppendToSlice(slice2, ch2)

// 	// fmt.Println(slice1)
// 	// fmt.Println(slice2)
// 	return Equal(slice1, slice2)

// }

func AppendToSlice(slice []int, ch chan int) []int {
	for i := range ch {
		slice = append(slice, i)
	}
	return slice
}

func Equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
