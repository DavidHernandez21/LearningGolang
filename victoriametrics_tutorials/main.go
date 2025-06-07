// Go program that demonstrates graceful shutdown of an HTTP server
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

func main() {
	// Setup signal context
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Readiness endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if isShuttingDown.Load() {
			http.Error(w, "Shutting down", http.StatusServiceUnavailable)
			return
		}

		_, err := fmt.Fprintln(w, "OK")
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}
	})

	// Sample business logic
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
			_, err := fmt.Fprintln(w, "Hello, world!")
			if err != nil {
				log.Printf("Error writing response: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}
		case <-r.Context().Done():
			http.Error(w, "Request cancelled.", http.StatusRequestTimeout)
		}
	})

	// Ensure in-flight requests aren't cancelled immediately on SIGTERM
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	server := &http.Server{
		Addr: ":8080",
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Server starting on :8080.")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for signal
	<-rootCtx.Done()
	stop()
	isShuttingDown.Store(true)
	log.Println("Received shutdown signal, shutting down.")

	// Give time for readiness check to propagate
	time.Sleep(_readinessDrainDelay)
	log.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()

	err := server.Shutdown(shutdownCtx)

	stopOngoingGracefully()

	if err != nil {
		log.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	log.Println("Server shut down gracefully.")
}
