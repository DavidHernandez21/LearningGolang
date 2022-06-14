package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	// _ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"unicode"
)

const (
	cpuProfile = "cpu.prof"
	memProfile = "mem.prof"
)

var byteArray [1]byte

func readByte(r io.Reader) (rune, error) {
	// var b [1]byte
	_, err := r.Read(byteArray[:])
	return rune(byteArray[0]), err
}

func main() {
	// f1, err := os.Create(cpuProfile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pprof.StartCPUProfile(f1)
	// defer f1.Close()
	// defer pprof.StopCPUProfile()
	f1, err := os.Create(memProfile)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f1.Close() // error handling omitted for example
	// runtime.MemProfile()
	runtime.MemProfileRate = 1
	// runtime.GC() // get up-to-date statistics

	f, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatalf("failed to open file %q: %v", os.Args[1], err)
	}
	defer f.Close()

	words := 0
	inword := false
	buffer := bufio.NewReader(f)
	// var byteArray [1]byte
	// byteArray := [1]byte{}

	for {
		r, err := readByte(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read byte: %v", err)
		}
		if unicode.IsSpace(r) && inword {
			words++
			inword = false
		}
		inword = unicode.IsLetter(r)

	}
	fmt.Printf("%q: %d words\n", os.Args[1], words)

	if err := pprof.Lookup("heap").WriteTo(f1, 0); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

}
