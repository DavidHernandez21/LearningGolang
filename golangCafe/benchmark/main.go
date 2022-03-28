package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	cpuProfile = "cpu.prof"
	memprofile = "mem.prof"
)

func main() {

	var m = make(map[int]int)

	mytimer := StartTimer("main")
	defer mytimer()

	// f1, err := os.Create(cpuProfile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pprof.StartCPUProfile(f1)
	// defer closeFile(f1)
	// defer pprof.StopCPUProfile()

	fiboQuick := quickFibonacciNoVarDeclaration(100, m)
	log.Println(fiboQuick)

	f, err := os.Create(memprofile)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer closeFile(f) // error handling omitted for example
	runtime.GC()       // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

}

func quickFibonacci(n int, m map[int]int) int {

	if n < 2 {
		return n
	}

	var f int
	if v, ok := m[n]; ok {
		f = v
	} else {
		f = quickFibonacci(n-2, m) + quickFibonacci(n-1, m)
		m[n] = f
	}

	return f
}

func quickFibonacciNoVarDeclaration(n int, m map[int]int) int {

	if n < 2 {
		return n
	}

	if v, ok := m[n]; ok {

		return v
	}

	var f int = quickFibonacciNoVarDeclaration(n-2, m) + quickFibonacciNoVarDeclaration(n-1, m)
	m[n] = f

	return f
}

func fibonacci(n int) int {
	if n < 2 {
		return n
	}

	return fibonacci(n-1) + fibonacci(n-2)
}

func StartTimer(name string) func() {
	t := time.Now()
	log.Println(name, "started")

	return func() {
		// d := time.Now().Sub(t)
		d := time.Since(t)
		log.Println(name, "took", d)
	}

}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Fatalf("error closing file '%v': %v", file.Name(), err)
	}
}
