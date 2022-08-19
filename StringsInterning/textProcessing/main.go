package main

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"unsafe"
)

// const (
// 	memProfile = "mem.prof"
// )

type stringInterner map[string]string

func newStringInterner(c int) stringInterner {
	return make(stringInterner, c)
}

func (si stringInterner) InternBytes(b []byte) string {
	// s := string(b)
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

// func linesBytesCount(s []byte) int {
// 	nl := []byte{'\n'}
// 	n := bytes.Count(s, nl)
// 	if len(s) > 0 && !bytes.HasSuffix(s, nl) {
// 		n++
// 	}
// 	return n
// }

func main() {
	// f1, err := os.Create(memProfile)
	// if err != nil {
	// 	log.Fatal("could not create memory profile: ", err)
	// }
	// defer f1.Close() // should be a noop

	f, err := os.Open("1984.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	// f1, _ := os.ReadFile("1984.txt")
	// log.Println("bytes count", linesBytesCount(f1))

	st, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	si := newStringInterner(int(st.Size() / 4))
	//
	// si := stringInterner{}
	words := make([]string, 0, st.Size()/4)

	for scanner.Scan() {
		words = append(words, si.InternBytes(scanner.Bytes())) // intern words
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("file size div by 4", st.Size()/4)
	log.Println("len of words", len(words))
	log.Println(words[130:150])
	for _, word := range words[130:150] {
		log.Printf("%x ", stringptr(word))
	}

	// if err := pprof.Lookup("heap").WriteTo(f1, 0); err != nil {
	// 	log.Fatal("could not write memory profile: ", err)
	// }

	// err = f1.Close()
	// if err != nil {
	// 	log.Fatal("could not close memory profile: ", err)
	// }

	// p := [10]int{10, 25, 3, 46, 5, 6, 7, 8, 9, 10}

	// intsAppend := make([]int, 0, len(p))
	// instAssing := make([]int, len(p))
	// for i := range p {
	// 	intsAppend = append(intsAppend, p[i])
	// 	instAssing[i] = p[i]
	// }

	// log.Println("intsAppend", intsAppend)
	// log.Println("instAssing", instAssing)
}
