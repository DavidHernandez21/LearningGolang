package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/dependencyInjection/cache"
	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/dependencyInjection/database"
	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/dependencyInjection/handlers"

	"github.com/gorilla/mux"
)

func main() {

	var numTables uint64 = 5
	var numObjects uint64 = 10

	db := database.NewInMemoryDB(numTables, numObjects)
	cache := cache.NewInMemoryCache()

	repositories := handlers.NewDefaultRepositoryProvider(db, cache)

	server := handlers.NewApiServer(repositories)

	// router := http.NewServeMux()
	router := mux.NewRouter()
	router.Use(server.CreateProviderMiddleware)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/getAll/{table:[a-zA-Z]+}", handlers.GetAllRecords)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/insert/{table:[a-zA-Z]+}", handlers.InsertSample)
	postRouter.HandleFunc("/createTable/{table:[a-zA-Z]+}", handlers.CreateTable)

	logger := log.New(os.Stdout, "", log.LstdFlags)

	s := http.Server{
		Addr:         "localhost:8080",  // configure the bind address
		Handler:      router,            // set the default handler
		ErrorLog:     logger,            // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		// signal.Notify(c, os.Kill)

		// Block until a signal is received.
		sig := <-signalChan
		logger.Println("Got signal:", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			logger.Fatalf("error shutting down the server: %v", err)
		}

	}()

	log.Println("listening on port 8080")
	log.Fatal(s.ListenAndServe())

}
