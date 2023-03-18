package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	logger := log.New(os.Stdout, "memray ", log.LstdFlags)

	router := mux.NewRouter()

	getRouter := router.Methods(http.MethodGet).Subrouter()

	getRouter.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	getRouter.Handle("/1", http.StripPrefix("/1", http.FileServer(http.Dir("./static1"))))

	// create a new server
	server := http.Server{
		Addr:         "localhost:8080",  // configure the bind address
		Handler:      router,            // set the default handler
		ErrorLog:     logger,            // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	CtrlCHandler(logger, &server)

	logger.Println("Listening to port 8080")

	logger.Fatal(server.ListenAndServe())

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
