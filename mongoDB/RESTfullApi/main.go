package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/mongoDB/RESTfullApi/clients"

	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/mongoDB/RESTfullApi/handlers"
)

// var client *mongo.Client

func main() {

	logger := log.New(os.Stdout, "mongoDBAtlas-api ", log.LstdFlags)

	client := clients.ConnectClient(logger)

	collection := client.Database("thepolyglotdeveloper").Collection("people")

	// create the handlers
	GetPersonByIdEndpoint := handlers.NewGetPersonByIdEndpoint(logger, collection)
	CreatePersonEndpoint := handlers.NewCreatePersonEndpoint(logger, collection)
	GetPeopleEndpoint := handlers.NewGetPeopleEndpoint(logger, collection)
	GetPersonByNameEndpoint := handlers.NewGetPersonByNameEndpoint(logger, collection)

	logger.Println("Starting the application...")

	clients.CtrlCHandler(client, logger)

	defer func() {
		clients.DisconnectClient(client, logger)
	}()

	router := mux.NewRouter()

	getRouter := router.Methods(http.MethodGet).Subrouter()

	getRouter.HandleFunc("/person/{id}", GetPersonByIdEndpoint.ServeHTTP)
	getRouter.HandleFunc("/people", GetPeopleEndpoint.ServeHTTP)
	getRouter.HandleFunc("/personName/{name}", GetPersonByNameEndpoint.ServeHTTP)

	// router.HandleFunc("/person", CreatePersonEndpoint.ServeHTTP).Methods("POST")
	// router.HandleFunc("/people", GetPeopleEndpoint.ServeHTTP).Methods("GET")
	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/person", CreatePersonEndpoint.ServeHTTP)
	// router.HandleFunc("/person/{id}", GetPersonEndpoint.ServeHTTP).Methods("GET")
	// router.HandleFunc("/test", test).Methods("GET")
	logger.Fatal(http.ListenAndServe("localhost:8080", router))

}
