// package main

// import (
// 	// "fmt"
// 	"log"
// 	"math/rand"
// 	"os"
// 	"time"
// )

// func main() {

// 	rand.Seed(time.Now().UnixNano())

// 	max := 20

// 	min := 10

// 	wait := rand.Intn(max-min) + min

// 	ch := make(chan struct{})

// 	// quick source. could be from some external request
// 	go externalMessagesNonBlocking(ch)

// 	// slow receiver
// 	for {
// 		select {
// 		case <-ch:
// 			log.Println("received message")
// 			// slow operation
// 			go func() {
// 				time.Sleep(time.Second * 3)
// 				log.Println("slow operation ended")
// 			}()
// 		case <-time.After(time.Duration(wait) * time.Millisecond):
// 			log.Println("timeout")
// 			os.Exit(1)
// 		}
// 	}
// }

// func externalMessagesNonBlocking(ch chan<- struct{}) {
// 	for {
// 		go func() {
// 			log.Println("created goroutine")
// 			select {
// 			case ch <- struct{}{}: // successfully send the message
// 				log.Println("send message")
// 			default: // or <-timer.C for timeout
// 				log.Println("Default case, receiver was blocked")
// 				// receiver was blocked, do something
// 			}
// 		}()
// 		time.Sleep(time.Second)
// 	}
// }
