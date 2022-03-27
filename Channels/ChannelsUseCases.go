// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
// )

// // func longTimeRequest() <-chan int32 {
// // 	r := make(chan int32)

// // 	go func() {
// // 		// Simulate a workload.
// // 		time.Sleep(time.Second * 3)
// // 		r <- rand.Int31n(100)
// // 	}()

// // 	return r
// // }

// // func sumSquares(a, b int32) int32 {
// // 	return a*a + b*b
// // }

// // func main() {
// // 	rand.Seed(time.Now().UnixNano())

// // 	a, b := longTimeRequest(), longTimeRequest()

// // 	m, n := <-a, <-b
// // 	fmt.Println(m, n)
// // 	fmt.Println(sumSquares(m, n))
// // }

// func longTimeRequest(r chan<- int32) {
// 	// Simulate a workload.
// 	time.Sleep(time.Second * 3)
// 	r <- rand.Int31n(100)
// }

// func sumSquares(a, b int32) int32 {
// 	return a*a + b*b
// }

// func main() {
// 	rand.Seed(time.Now().UnixNano())

// 	results := make(chan int32)
// 	go longTimeRequest(results)
// 	go longTimeRequest(results)

// 	fmt.Println(sumSquares(<-results, <-results))
// }
