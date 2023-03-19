package main

import (
	"fmt"
	"log"
	"src/github.com/DavidHernandez21/MemStoreProject/io"
	"src/github.com/DavidHernandez21/MemStoreProject/ms"
)

type (
	// ReadWriter combines the interfaces io.Reader and io.Writer
	ReadWriter interface {
		io.Reader
		io.Writer
	}
)

func main() {
	var m, _ = ms.NewMemStore(100)
	// test(m)

	errChan := make(chan error)
	totalBytesChan := make(chan int)
	readChan := make(chan []byte)

	go writeUntilError(m, errChan, totalBytesChan)

	go readAfterMsIsFull(m, errChan, readChan)

	// fmt.Println(<-errChan)
	fmt.Println(<-totalBytesChan)

	fmt.Println(string(<-readChan))

}

// func test(rw ReadWriter) {
// 	var totalBytes int
// 	n, _ := rw.Write([]byte("Hello, world"))
// 	totalBytes += n

// 	n, _ = rw.Write([]byte(". "))
// 	totalBytes += n

// 	n, _ = rw.Write([]byte("It is a wonderful day."))
// 	totalBytes += n

// 	b := make([]byte, 1)
// 	_, err := rw.Read(b)
// 	for err == nil {
// 		fmt.Print(string(b))
// 		_, err = rw.Read(b)
// 	}
// 	fmt.Printf("\n%v\n", err) // should print "Hello, world. It is a wonderful day."
// }

func writeUntilError(rw ReadWriter, errChan chan error, totalBytesChan chan int) {

	var totalBytes int
	n, err := rw.Write([]byte("Hello, world "))
	totalBytes += n

	for err == nil {
		n, err = rw.Write([]byte("Hello, world "))
		totalBytes += n

	}

	errChan <- err
	totalBytesChan <- totalBytes

}

func readAfterMsIsFull(rw ReadWriter, errChan chan error, readChan chan []byte) {

	<-errChan

	b := make([]byte, 101)
	_, err := rw.Read(b)

	if err != nil {
		log.Printf("error while reading from the buffer: %v", err)
	}

	readChan <- b

}
