package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func hashAndSend(r io.Reader, lenght int) (string, []byte) {
	w := sha256.New()

	//any reads from tee will read from r and write to w
	tee := io.TeeReader(r, w)

	buff, err := sendReader(tee, lenght)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return hex.EncodeToString(w.Sum(nil)), buff
	// fmt.Println(sha)
}

// sendReader sends the contents of an io.Reader to stdout using a 256 byte buffer
func sendReader(data io.Reader, lenght int) ([]byte, error) {
	buff := make([]byte, lenght)
	for {
		_, err := data.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			// fmt.Println(err)
			return nil, err
		}

	}
	// fmt.Println(len(buff))
	// fmt.Print(strings.TrimSuffix(string(buff), " "))
	// fmt.Println("")
	return buff, nil
}

func main() {
	repeat := 300
	r1 := strings.NewReader(strings.Repeat("Daje Roma ", repeat)) //our io.Reader
	// r2 := strings.NewReader("hello world")

	// hashAndSendNaive(r1)
	// fmt.Println(r1.Size())
	// fmt.Println(r1.Len(), r1.Size())
	shaString, buff := hashAndSend(r1, r1.Len())
	fmt.Println(strings.TrimSuffix(string(buff), " "))
	fmt.Println(shaString)
}
