// Description: This code demonstrates the use of a buffered channel with select statements to attempt sending and receiving messages.
package main

import "fmt"

func main4() {
	type Book struct {
		id int
	}

	bookshelf := make(chan Book, 3)

	for i := range cap(bookshelf) * 2 {
		select {
		case bookshelf <- Book{id: i}:
			fmt.Println("succeeded to put book", i)
		default:
			fmt.Println("failed to put book")
		}
	}

	for range cap(bookshelf) * 2 {
		select {
		case book := <-bookshelf:
			fmt.Println("succeeded to get book", book.id)
		default:
			fmt.Println("failed to get book")
		}
	}
}
