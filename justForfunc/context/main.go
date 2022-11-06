package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	mySleepAndTalk(ctx, 5*time.Second, "hello")
}

func mySleepAndTalk(ctx context.Context, d time.Duration, msg string) {
	delay := time.NewTimer(d)
	select {
	case <-delay.C:
		fmt.Println(msg)
	case <-ctx.Done():
		log.Printf("Context error: %v", ctx.Err())
		if !delay.Stop() {
			delay.C = nil
		}
		// default:
		// 	log.Println("default case")
	}
}
