package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
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
			logrus.Error("Error: ", err)
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
	logrus.Infof("request to external at %v\n", time.Now())
	d := time.Duration(rg.Int31n(500)) * time.Millisecond
	time.Sleep(d)
}
