package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sync/semaphore"
)

var host string
var ports string
var timeout int

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host to scan.")
	flag.StringVar(&ports, "ports", "5400-5500", "Port(s) (e.g. 80, 22-100).")
	flag.IntVar(&timeout, "timeout", 10, "Timeout in seconds (default is 5).")

	rand.Seed(time.Now().UnixNano())
}

func main() {

	timer := StartTimer("main")
	defer timer()

	flag.Parse()

	portsToScan, err := parsePortsToScan(ports)
	if err != nil {
		fmt.Printf("Failed to parse ports to scan: %s\n", err)
		os.Exit(1)
	}

	openPorts := make([]int, 0, len(portsToScan))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		printResults(openPorts)
		os.Exit(0)
	}()

	var semMaxWeight int64 = 100_000
	var semAcquisitionWeight int64 = 100

	sem := semaphore.NewWeighted(semMaxWeight)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	for k := range portsToScan {
		if err := sem.Acquire(ctx, semAcquisitionWeight); err != nil {
			fmt.Printf("Failed to acquire semaphore (port %d): %v\n", portsToScan[k], err)
			break
		}

		go func(port int) {
			defer sem.Release(semAcquisitionWeight)
			fmt.Printf("num of active goroutines: %v\n", runtime.NumGoroutine())
			sleepy(1)
			p := scan(host, port)
			if p != 0 {
				openPorts = append(openPorts, p)
			}
		}(portsToScan[k])
	}

	// We block here until done.
	err = sem.Acquire(ctx, semMaxWeight)
	if err != nil {
		log.Fatalf("failed to acquire semaphore: %v", err)
	}

	printResults(openPorts)
}

func parsePortsToScan(portsFlag string) ([]int, error) {
	p, err := strconv.Atoi(portsFlag)
	if err == nil {
		return []int{p}, nil
	}

	ports := strings.Split(portsFlag, "-")
	if len(ports) != 2 {
		return nil, errors.New("unable to determine port(s) to scan")
	}

	minPort, err := strconv.Atoi(ports[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to a valid port number", ports[0])
	}

	maxPort, err := strconv.Atoi(ports[1])
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to a valid port number", ports[1])
	}

	if minPort <= 0 || maxPort <= 0 {
		return nil, fmt.Errorf("port numbers must be greater than 0")
	}

	// var results []int
	results := make([]int, 0, maxPort-minPort+1)
	for p := minPort; p <= maxPort; p++ {
		results = append(results, p)
	}
	return results, nil
}

func scan(host string, port int) int {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("%d CLOSED (%s)\n", port, err)
		return 0
	}
	conn.Close()
	return port
}

func sleepy(max int) {
	n := rand.Intn(max)
	time.Sleep(time.Duration(n) * time.Second)
}

func printResults(ports []int) {
	sort.Ints(ports)
	fmt.Println("\nResults\n--------------")
	for k := range ports {
		fmt.Printf("%d - open\n", ports[k])
	}
}

func StartTimer(name string) func() {
	t := time.Now()
	fmt.Println(name, "started")

	return func() {
		// d := time.Now().Sub(t)
		d := time.Since(t)
		fmt.Println(name, "took", d)
	}

}
