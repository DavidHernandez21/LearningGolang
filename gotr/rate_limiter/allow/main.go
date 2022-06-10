package main

import (
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

	for i := 0; i < 10; {
		if lim.Allow() {
			i++
			callExternal()
			continue
		}
		// logrus.Infof("not allowed to call yet: %v", time.Now())
		// d := time.Duration(rg.Int31n(500)) * time.Millisecond
		// time.Sleep(time.Duration(rg.Int31n(500)) * time.Millisecond)

	}
}

func callExternal() {
	logrus.Infof("request to external at %v\n", time.Now())
	// d := time.Duration(rg.Int31n(500)) * time.Millisecond
	time.Sleep(time.Duration(rg.Int31n(500)) * time.Millisecond)
}
