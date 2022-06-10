package main

import (
	"context"
	"math/rand"
	"time"

	log "github.com/mgutz/logxi/v1"
	"golang.org/x/time/rate"
)

var (
	rg = rand.New(rand.NewSource(time.Now().Unix()))
)

func main() {
	r := rate.Every(2 * time.Second)
	lim := rate.NewLimiter(r, 3)
	n := lim.Burst()
	counter := 0

	for i := 0; i < 10; i++ {

		err := lim.WaitN(context.Background(), n)
		if err != nil {
			log.Error("Error: ", err)
		}
		for j := 0; j < n; j++ {
			if counter == 10 {
				return
			}
			callExternal()
			counter++
		}
	}
}

func callExternal() {
	log.Info("request to external at %v\n", time.Now())
	d := time.Duration(rg.Int31n(500)) * time.Millisecond
	time.Sleep(d)
}
