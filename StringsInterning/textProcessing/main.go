package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"
)

type stringInterner map[string]string

func (si stringInterner) InternBytes(b []byte) string {
	if interned, ok := si[string(b)]; ok {
		return interned
	}
	s := string(b)
	si[s] = s
	return s
}

// stringptr returns a pointer to the string data.
func stringptr(s string) uintptr {
	return (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
}

func main() {
	f, err := os.Open("1984.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	si := stringInterner{}
	var words []string
	for scanner.Scan() {
		words = append(words, si.InternBytes(scanner.Bytes())) // intern words
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// log.Print(len(words))
	log.Println(words[111:122])
	for _, word := range words[111:122] {
		fmt.Printf("%x ", stringptr(word))
	}
}
