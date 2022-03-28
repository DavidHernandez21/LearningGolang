package main

import (
	"bufio"
	"log"
	"os"
)

type (
	stringInternerStruct struct {
		mapper map[string]string
	}

	stringInterner map[string]string
)

func newStringInterner(cap int) *stringInternerStruct {
	return &stringInternerStruct{
		mapper: make(map[string]string, cap),
	}
}

func (si *stringInternerStruct) InternBytes(b []byte) string {
	if interned, ok := si.mapper[string(b)]; ok {
		return interned
	}
	s := string(b)
	si.mapper[s] = s
	return s
}

func (si stringInterner) InternBytes(b []byte) string {
	if interned, ok := si[string(b)]; ok {
		return interned
	}
	s := string(b)
	si[s] = s
	return s
}

func main() {

	f, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	si := newStringInterner(1000)
	var words []string
	for scanner.Scan() {
		words = append(words, si.InternBytes(scanner.Bytes())) // intern words
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println(len(words))

}
