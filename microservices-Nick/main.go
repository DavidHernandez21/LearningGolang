package main

import (
	"context"
	"src/github.com/DavidHernandez21/microservices-Nick/handlers"
	"syscall"

	// "github.com/nicholasjackson/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

// var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	// env.Parse()

	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	if err := godotenv.Load(".env"); err != nil {
		l.Println("No .env file found")
	}

	bindAddress := os.Getenv("bindAddress")

	if bindAddress == "" {
		bindAddress = "localhost:8080"
	}

	// l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	// create the handlers
	ph := handlers.NewProducts(l)

	// create a new serve mux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", ph)

	// create a new server
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	errShutdown := s.Shutdown(ctx)

	if errShutdown != nil {
		l.Println("Error while trying to gracefully shutdown server")
		os.Exit(1)
	}
}
