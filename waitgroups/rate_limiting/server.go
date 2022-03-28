package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	li, err := net.Listen("tcp", ":8080")

	defer closeListener(li)

	stopSignalHandler(li)

	if err != nil {
		log.Fatalf("could not create listener: %v", err)
	}

	log.Println("tcp server started, listennig to port 8080")

	var connections int32
	for {
		conn, err := li.Accept()
		if err != nil {
			continue
		}
		connections++

		go func(conn net.Conn) {
			defer func() {
				_ = conn.Close()
				atomic.AddInt32(&connections, -1)
			}()
			if atomic.LoadInt32(&connections) > 3 {
				_, err := conn.Write([]byte("too many connections"))
				if err != nil {
					log.Fatalf("could not write to connection: %v", err)
				}
				log.Println("concurrent connections over threshold of 3")
				return
			}

			// simulate heavy work
			time.Sleep(time.Second)
			_, err := conn.Write([]byte("success"))
			if err != nil {
				log.Fatalf("could not write to connection: %v", err)
			}
		}(conn)
	}
}

func closeListener(lis net.Listener) {
	err := lis.Close()
	if err != nil {
		log.Fatalf("Could not close listener: %v", err)
	}

}

func stopSignalHandler(lis net.Listener) {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		log.Printf("- Ctrl+C pressed, exiting\n Signal recieved: %v\n", sig)
		closeListener(lis)
		os.Exit(0)
	}()

}
