package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

var (
	defaultServer = "http://localhost:8080"
)

const (
	memProfile = "mem.prof"
)

func main() {
	timeit := StartTimer("client")
	defer timeit()

	f1, err := os.Create(memProfile)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f1.Close() // should be a noop

	var url string
	flag.StringVar(&url, "s", defaultServer, "server in the form http://host:port")
	flag.Parse()

	submitRequests(url)

	if err := pprof.Lookup("heap").WriteTo(f1, 0); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	err = f1.Close()
	if err != nil {
		log.Fatal("could not close memory profile: ", err)
	}
}
