package main

import (
	"context"
	"fmt"
	"time"
)

type requestID string

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	ctx = context.WithValue(ctx, requestID("requestID"), "12345")
	for i := range gen(ctx) {
		time.Sleep(1 * time.Millisecond)
		fmt.Println("Got from gen:", i)
		if i == 10 {
			break
		}
	}
	cancel()
	time.Sleep(1 * time.Second)
}

func gen(ctx context.Context) <-chan int {
	out := make(chan int)
	var i int
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Request ID:", ctx.Value(requestID("requestID")))
				fmt.Println("stopping goroutine...")
				fmt.Println("error: ", ctx.Err())
				close(out)
				return
			case out <- i:
				i++
			}
		}
	}()
	return out
}
