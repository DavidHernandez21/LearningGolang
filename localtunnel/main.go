package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	localtunnel "github.com/jonasfj/go-localtunnel"
)

func main() {

	// Setup a listener for localtunnel
	listener, err := localtunnel.Listen(localtunnel.Options{})

	if err != nil {
		log.Fatalf("error localtunnel listener: %v", err)
	}
	logger := log.New(os.Stdout, "localtunnel ", log.LstdFlags)

	server := http.Server{
		Addr:         "localhost:8080",  // configure the bind address
		Handler:      nil,               // set the default handler
		ErrorLog:     logger,            // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Handle request from localtunnel
	CtrlCHandler(logger, &server)

	logger.Printf("locatunnel URL: %v", listener.URL())

	logger.Fatal(server.Serve(listener))

}

func CtrlCHandler(logger *log.Logger, server *http.Server) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		logger.Printf("Sutting down server. Signal recieved: %v\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("error shutting down the server: %v", err)
		}
		logger.Println("server shut down")
		os.Exit(0)
	}()
}
