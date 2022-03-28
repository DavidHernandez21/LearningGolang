// package main

// import (
// 	"fmt"
// 	"log"
// 	"sync"
// 	"time"
// )

// type item struct {
// 	price    float32
// 	category string
// }

// func gen(items ...item) <-chan item {
// 	out := make(chan item, len(items))
// 	for _, i := range items {
// 		out <- i
// 	}
// 	close(out)
// 	return out
// }

// func discount(items <-chan item) <-chan item {
// 	out := make(chan item)
// 	go func() {
// 		defer close(out)
// 		for i := range items {
// 			time.Sleep(time.Second / 2)
// 			if i.category == "shoe" {
// 				i.price /= 2
// 			}
// 			out <- i
// 			// select {
// 			// case out <- i:
// 			// case <-done:
// 			// 	fmt.Println("we are done with discount(s)")
// 			// 	return
// 			// }

// 		}

// 	}()
// 	return out
// }

// func fanIn(channelmap map[<-chan item]int, channels ...<-chan item) <-chan item {
// 	var wg sync.WaitGroup

// 	out := make(chan item)

// 	output := func(c <-chan item) {
// 		defer wg.Done()
// 		for i := range c {
// 			out <- i
// 			fmt.Printf("fan in recieved value: %+v from channel: %v\n", i, channelmap[c])
// 			// select {
// 			// case out <- i:
// 			// 	fmt.Printf("fan in recieved value: %+v from channel: %v\n", i, channelmap[c])
// 			// case <-done:
// 			// 	fmt.Println("we are done fanning in")
// 			// 	return
// 			// }
// 		}
// 	}

// 	wg.Add(len(channels))

// 	for _, c := range channels {

// 		go output(c)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(out)
// 	}()
// 	return out

// }

// func timeTrack(start time.Time, name string, workers int16) {
// 	elapsed := time.Since(start)
// 	log.Printf("%s took %s, with %v worker(s)", name, elapsed, workers)
// }

// func main() {
// 	// done := make(chan struct{})
// 	// defer close(done)
// 	var workers int16
// 	workers = 7
// 	defer timeTrack(time.Now(), "main", workers)

// 	c := gen(
// 		item{price: 8, category: "shirt"},
// 		item{price: 20, category: "shoe"},
// 		item{price: 24, category: "shoe"},
// 		item{price: 4, category: "drink"},
// 		item{price: 9, category: "ball"},
// 		item{price: 16, category: "pants"},
// 		item{price: 17, category: "shoe"})

// 	channels := make([]<-chan item, 0, workers)
// 	channelsMap := make(map[<-chan item]int)

// 	for i := 0; i < int(workers); i++ {
// 		c := discount(c)
// 		channels = append(channels, c)
// 		channelsMap[c] = i + 1
// 		fmt.Printf("address of channel %v: %v\n", i+1, c)
// 	}

// 	// c1 := discount(c)
// 	// c2 := discount(c)
// 	// c3 := discount(c)
// 	// fmt.Printf("channel 1 address: %v, channel 2 address: %v\n", channels)

// 	// fmt.Println(<-out)
// 	out := fanIn(channelsMap, channels...)
// 	for processed := range out {
// 		fmt.Printf("Category: %v\t Price: %v\n", processed.category, processed.price)
// 	}

// }
