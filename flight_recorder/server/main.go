// Package main implements a server that records flight recorder snapshots
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/trace"
	"strconv"
	"sync"
	"time"
)

// bucket is a simple mutex-protected counter.
type bucket struct {
	mu      sync.Mutex
	guesses int
}

func main() {
	// Set up the flight recorder
	fr := trace.NewFlightRecorder(trace.FlightRecorderConfig{
		MinAge:   200 * time.Millisecond,
		MaxBytes: 1 << 20, // 1 MiB
	})
	if err := fr.Start(); err != nil {
		log.Fatalf("failed to start flight recorder: %s", err)
	}

	// Make one bucket for each valid number a client could guess.
	// The HTTP handler will look up the guessed number in buckets by
	// using the number as an index into the slice.
	buckets := make([]bucket, 100)

	// Every minute, we send a report of how many times each number was guessed.
	go func() {
		for range time.Tick(1 * time.Minute) {
			sendReport(buckets)
		}
	}()

	// Choose the number to be guessed.
	answer := rand.Intn(len(buckets))

	http.HandleFunc("/guess-number", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Fetch the number from the URL query variable "guess" and convert it
		// to an integer. Then, validate it.
		guess, err := strconv.Atoi(r.URL.Query().Get("guess"))
		if err != nil || (0 > guess || guess >= len(buckets)) {
			http.Error(w, "invalid 'guess' value", http.StatusBadRequest)
			return
		}

		// Select the appropriate bucket and safely increment its value.
		b := &buckets[guess]
		b.mu.Lock()
		b.guesses++
		b.mu.Unlock()

		// Respond to the client with the guess and whether it was correct.
		if _, err := fmt.Fprintf(w, "guess: %d, correct: %t", guess, guess == answer); err != nil {
			log.Printf("failed to write response: %s", err)
		}

		// Capture a snapshot if the response takes more than 100ms.
		// Only the first call has any effect.
		if fr.Enabled() && time.Since(start) > 100*time.Millisecond {
			log.Printf("capturing snapshot for slow request: duration=%s", time.Since(start))
			go captureSnapshot(fr)
		}

		log.Printf("HTTP request: endpoint=/guess-number guess=%d duration=%s", guess, time.Since(start))
	})
	server := &http.Server{
		Addr:              ":8090",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

// sendReport posts the current state of buckets to a remote service.
func sendReport(buckets []bucket) {
	counts := make([]int, len(buckets))

	for index := range buckets {
		b := &buckets[index]
		b.mu.Lock()

		counts[index] = b.guesses
		b.mu.Unlock()
	}

	// Marshal the report data into a JSON payload.
	b, err := json.Marshal(counts)
	if err != nil {
		log.Printf("failed to marshal report data: error=%s", err)
		return
	}
	url := "http://127.0.0.1:8091/guess-number-report"
	if resp, err := http.Post(url, "application/json", bytes.NewReader(b)); err != nil {
		log.Printf("failed to send report: %s", err)
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %s", err)
		}
	} else {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %s", err)
		}
	}
}

var once sync.Once

// captureSnapshot captures a flight recorder snapshot.
func captureSnapshot(fr *trace.FlightRecorder) {
	// once.Do ensures that the provided function is executed only once.
	once.Do(func() {
		f, err := os.Create("snapshot.trace")
		if err != nil {
			log.Printf("opening snapshot file %s failed: %s", f.Name(), err)
			return
		}
		defer func() {
			if err = f.Close(); err != nil {
				log.Printf("closing snapshot file %s failed: %s", f.Name(), err)
			}
		}()

		// WriteTo writes the flight recorder data to the provided io.Writer.
		_, err = fr.WriteTo(f)
		if err != nil {
			log.Printf("writing snapshot to file %s failed: %s", f.Name(), err)
			return
		}

		// Stop the flight recorder after the snapshot has been taken.
		fr.Stop()
		log.Printf("captured a flight recorder snapshot to %s", f.Name())
	})
}
