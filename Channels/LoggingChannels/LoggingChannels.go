package main

import (
	"fmt"
	"os"
	"time"
)

const (
	logInfo    = "INFO"
	logWarning = "WARNING"
	logError   = "ERROR"
)

type logEntry struct {
	time     time.Time
	severity string
	message  string
}

func main() {

	logCh := make(chan logEntry)

	doneCh := make(chan struct{})

	go logger(logCh)

	go selecLogger(logCh, doneCh)

	log1 := logEntry{
		time:     time.Now(),
		severity: logInfo,
		message:  "App is starting",
	}

	log2 := logEntry{
		time:     time.Now(),
		severity: logInfo,
		message:  "App is shutting down",
	}

	for i := 0; i < 5; i++ {
		go sendLog(log1, logCh)
		go sendLog(log2, logCh)
	}

	time.Sleep(time.Millisecond * 100)
	doneCh <- struct{}{}
	time.Sleep(time.Millisecond * 200)

}

func logger(ch chan logEntry) {
	for entry := range ch {
		fmt.Printf("%v - [%v]%v\n", entry.time.Format(time.RFC3339),
			entry.severity, entry.message)

	}

}

func selecLogger(ch chan logEntry, doneCh chan struct{}) {

	delay := time.NewTimer(1000 * time.Millisecond)
outerloop:
	for {
		select {
		case entry := <-ch:
			fmt.Printf("%v - [%v]%v\n", entry.time.Format(time.RFC3339),
				entry.severity, entry.message)

		case <-delay.C:
			fmt.Println("timeout")
			os.Exit(1)

		case <-doneCh:
			fmt.Println("recieve from done channel")
			break outerloop

		}
	}

	if !delay.Stop() {
		delay.C = nil
		fmt.Println("Freeing up resources")
	}
}

func sendLog(log logEntry, ch chan logEntry) {
	ch <- log
}

// func main() {
// 	var ball = make(chan string)
// 	kickBall := func(playerName string) {
// 		for {
// 			fmt.Println(<-ball, "kicked the ball.")
// 			time.Sleep(time.Second)
// 			ball <- playerName
// 		}
// 	}
// 	go kickBall("John")
// 	go kickBall("Alice")
// 	go kickBall("Bob")
// 	go kickBall("Emily")
// 	ball <- "referee" // kick off
// 	// var c chan bool   // nil
// 	// <-c               // blocking here for ever
// 	// c <- true
// 	fmt.Println("receiving value from booleand channel")

// }
