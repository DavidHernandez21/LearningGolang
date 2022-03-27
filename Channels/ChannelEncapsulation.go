// package main

// import (
// 	"fmt"
// 	"os"
// 	"time"
// )

// var counter = func(n int) chan<- chan<- int {
// 	requests := make(chan chan<- int)
// 	go func() {
// 		for request := range requests {
// 			if request == nil {
// 				n++ // increase
// 			} else {
// 				fmt.Println("actual value of n: ", n)
// 				request <- n // take out
// 			}
// 		}
// 	}()

// 	// Implicitly converted to chan<- (chan<- int)

// 	return requests
// }(0)

// func main() {
// 	increase1000 := func(done chan<- struct{}) {
// 		for i := 0; i < 1000; i++ {
// 			counter <- nil
// 		}
// 		time.Sleep(2 * time.Second)
// 		done <- struct{}{}
// 	}

// 	done := make(chan struct{})

// 	for i := 0; i < 4; i++ {
// 		go increase1000(done)
// 		<-done
// 	}
// 	// go increase1000(done)
// 	// go increase1000(done)
// 	// <-done
// 	// <-done

// 	request := make(chan int)
// 	request2 := make(chan int)
// 	counter <- request
// 	fmt.Println(<-request) // 2000

// 	counter <- request2
// 	fmt.Println(<-request2)

// 	for {
// 		go increase1000(done)
// 		select {
// 		case <-done:
// 			counter <- request
// 			fmt.Println(<-request)

// 			// time.Sleep(2 * time.Second)
// 		case <-time.After(1 * time.Second):
// 			fmt.Println("timeout")
// 			os.Exit(1)
// 		}
// 	}

// }
