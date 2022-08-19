// Including main package
package main

// Importing fmt, io, and bytes
import (
	"bytes"
	"fmt"
	"io"
	"log"
)

// Calling main
func main() {

	// Calling Pipe method
	pipeReader, pipeWriter := io.Pipe()

	// Using Fprint in go function to write
	// data to the file
	go func() {
		fmt.Fprint(pipeWriter, "GeeksforGeeks\nis\na\nCS-Portal.\n")

		// Using Close method to close write
		pipeWriter.Close()
	}()

	// Creating a buffer
	buffer := new(bytes.Buffer)

	// Calling ReadFrom method and writing
	// data into buffer
	_, err := buffer.ReadFrom(pipeReader)
	if err != nil {
		log.Fatalln(err)
	}

	// Prints the data in buffer
	fmt.Print(buffer.String())
}
