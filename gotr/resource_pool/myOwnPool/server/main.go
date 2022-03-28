package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/mgutz/logxi/v1"
)

func main() {
	srv := newTCPServer("8080")
	err := srv.Start()
	if err != nil {
		log.Error("Failed to start TCPServer", err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Info("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	errShutdown := srv.Shutdown(ctx)

	if errShutdown != nil {
		log.Error("Error while trying to gracefully shutdown server:")
		os.Exit(1)
	}
	// d := 60 * time.Second
	// fmt.Printf("Sleeping for %v\n", d)
	// time.Sleep(d)
	// srv.Stop()
}
