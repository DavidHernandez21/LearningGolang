package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	log "github.com/mgutz/logxi/v1"
	"golang.org/x/time/rate"
)

var (
	rg = rand.New(rand.NewSource(time.Now().Unix()))
)

func main() {
	my_rate := rate.Every(2 * time.Second)
	lim := rate.NewLimiter(my_rate, 3)

	for i := 0; i < 10; i++ {
		err := lim.Wait(context.Background())
		if err != nil {
			log.Error("Error: ", err)
		}
		callExternal()
	}
}

func callExternal() {
	fmt.Printf("request to external at %v\n", time.Now())
	d := time.Duration(rg.Int31n(500)) * time.Millisecond
	time.Sleep(d)
}
