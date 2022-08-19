package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/mongoDB/RESTfullApi/data"
	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/timer"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// var client = clients.ConnectClient()

type (
	CreatePersonEndpoint struct {
		Logger     *log.Logger
		Collection *mongo.Collection
	}

	GetPersonByIdEndpoint struct {
		Logger     *log.Logger
		Collection *mongo.Collection
	}

	GetPeopleEndpoint struct {
		Logger     *log.Logger
		Collection *mongo.Collection
	}

	GetPersonByNameEndpoint struct {
		Logger     *log.Logger
		Collection *mongo.Collection
	}
)

func (c *CreatePersonEndpoint) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	stop := timer.StartTimer("CreatePersonEndpoint", c.Logger)

	defer stop()

	response.Header().Set("content-type", "application/json")
	var person data.Person
	// err := json.NewDecoder(request.Body).Decode(&person)
	err := person.FromJSON(request.Body)

	// c.Logger.Println(person)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		c.Logger.Printf("Error while marshalling the request body: %v\n", err)
		return
	}

	// collection := c.Client.Database("thepolyglotdeveloper").Collection("people")
	collection := c.Collection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err1 := collection.InsertOne(ctx, person)

	if err1 != nil {
		// response.WriteHeader(http.StatusInternalServerError)
		http.Error(response, "Internal Error", http.StatusInternalServerError)
		c.Logger.Printf("Error while inserting the person: %v \n%v\n", person, err1)
		return
	}

	err2 := json.NewEncoder(response).Encode(result)
	// err2 := person.ToJSON(response)

	if err2 != nil {
		// response.WriteHeader(http.StatusInternalServerError)
		http.Error(response, "Internal Error", http.StatusInternalServerError)
		c.Logger.Printf("Error while marshalling the result: %v \n%v\n", result, err2)
		return
	}

}

func (c *GetPersonByNameEndpoint) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	stop := timer.StartTimer("GetPersonByNameEndpoint", c.Logger)

	defer stop()

	response.Header().Set("content-type", "application/json")

	name := mux.Vars(request)["name"]

	// var person data.Person
	// collection := c.Client.Database("thepolyglotdeveloper").Collection("people")
	collection := c.Collection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	pattern := fmt.Sprintf("^%v.*", name)

	// c.Logger.Printf("patter: %v", pattern)

	regexValue := primitive.Regex{Pattern: pattern, Options: "i"}

	cursor, err := collection.Find(ctx, bson.D{primitive.E{Key: "firstname", Value: bson.D{primitive.E{Key: "$regex", Value: regexValue}}}})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, err := response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}

		return
	}

	defer cursor.Close(ctx)

	var people data.People

	people = appendPersonFromCursor(cursor, people, ctx, response, c.Logger)

	// for cursor.Next(ctx) {
	// 	// var person Person
	// 	var person data.Person
	// 	cursor.Decode(&person)
	// 	people = append(people, &person)
	// }
	// if err := cursor.Err(); err != nil {
	// 	response.WriteHeader(http.StatusInternalServerError)
	// 	_ , err := response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	// if err != nil {
	// 	c.Logger.Printf("Error while writing the error response: %v\n", err)
	// }
	// 	return
	// }

	// err1 := json.NewEncoder(response).Encode(people)
	err1 := people.ToJSON(response)

	if err1 == data.ErrNotFound {
		c.Logger.Printf("No Person was found with the name: %v", name)
		_, err := response.Write([]byte(`{ "message": "No Person was found with the name: '` + name + `'" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}
		return
	}

	if err1 != nil {

		response.WriteHeader(http.StatusInternalServerError)
		_, err := response.Write([]byte(`{ "message": "` + err1.Error() + `" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}
		return

	}

}

func (c *GetPersonByIdEndpoint) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	stop := timer.StartTimer("GetPersonByIdEndpoint", c.Logger)

	defer stop()

	response.Header().Set("content-type", "application/json")
	paramsId := mux.Vars(request)["id"]
	// paramsId := params["id"]
	id, errId := primitive.ObjectIDFromHex(paramsId)

	if errId != nil {
		response.WriteHeader(http.StatusInternalServerError)
		// http.Error(response, "Internal Error", http.StatusInternalServerError)
		_, err := response.Write([]byte(`{ "message": "` + errId.Error() + `" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}

		return
	}

	// log.Printf("params values: %v\n", params)

	// var person Person
	var person data.Person
	// collection := c.Client.Database("thepolyglotdeveloper").Collection("people")
	collection := c.Collection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err := collection.FindOne(ctx, data.Person{ID: id}).Decode(&person)

	if err == mongo.ErrNoDocuments {
		c.Logger.Printf("No Person was found with the id: %v", paramsId)
		_, err := response.Write([]byte(`{ "message": "No Person was found with the id: ` + paramsId + `" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}
		return
	}

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		// http.Error(response, "Internal Error", http.StatusInternalServerError)
		_, err := response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}
		return
	}
	// json.NewEncoder(response).Encode(person)
	err1 := person.ToJSON(response)

	if err1 != nil {

		response.WriteHeader(http.StatusInternalServerError)
		// http.Error(response, "Internal Error", http.StatusInternalServerError)
		_, err := response.Write([]byte(`{ "message": "` + err1.Error() + `" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}

		return

	}
}

func (c *GetPeopleEndpoint) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	stop := timer.StartTimer("GetPeopleEndpoint", c.Logger)

	defer stop()

	response.Header().Set("content-type", "application/json")
	// var people []Person
	var people data.People
	// collection := c.Client.Database("thepolyglotdeveloper").Collection("people")
	collection := c.Collection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, err := response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}

		return
	}
	defer cursor.Close(ctx)

	people = appendPersonFromCursor(cursor, people, ctx, response, c.Logger)

	// for cursor.Next(ctx) {
	// 	// var person Person
	// 	var person data.Person
	// 	cursor.Decode(&person)
	// 	people = append(people, &person)
	// }
	// if err1 := cursor.Err(); err1 != nil {
	// 	response.WriteHeader(http.StatusInternalServerError)
	// 	response.Write([]byte(`{ "message": "` + err1.Error() + `" }`))
	// 	return
	// }

	// err1 := json.NewEncoder(response).Encode(people)
	err2 := people.ToJSON(response)

	if err2 != nil {

		response.WriteHeader(http.StatusInternalServerError)
		_, err := response.Write([]byte(`{ "message": "` + err2.Error() + `" }`))
		if err != nil {
			c.Logger.Printf("Error while writing the error response: %v\n", err)
		}
		return

	}

}

func NewGetPersonByIdEndpoint(logger *log.Logger, collection *mongo.Collection) *GetPersonByIdEndpoint {
	return &GetPersonByIdEndpoint{
		Logger:     logger,
		Collection: collection,
	}
}

func NewCreatePersonEndpoint(logger *log.Logger, collection *mongo.Collection) *CreatePersonEndpoint {
	return &CreatePersonEndpoint{
		Logger:     logger,
		Collection: collection,
	}
}

func NewGetPeopleEndpoint(logger *log.Logger, collection *mongo.Collection) *GetPeopleEndpoint {
	return &GetPeopleEndpoint{
		Logger:     logger,
		Collection: collection,
	}
}

func NewGetPersonByNameEndpoint(logger *log.Logger, collection *mongo.Collection) *GetPersonByNameEndpoint {
	return &GetPersonByNameEndpoint{
		Logger:     logger,
		Collection: collection,
	}
}

func appendPersonFromCursor(cursor *mongo.Cursor, people data.People, ctx context.Context, response http.ResponseWriter, logger *log.Logger) data.People {

	stop := timer.StartTimer("appendPersonFromCursor", logger)

	defer stop()

	// var wg sync.WaitGroup

	for cursor.Next(ctx) {

		// wg.Add(1)
		// var person Person
		// go func() {
		// 	defer wg.Done()
		// 	var person data.Person
		// 	cursor.Decode(&person)
		// 	people = append(people, &person)
		// }()
		var person data.Person
		err := cursor.Decode(&person)
		if err != nil {
			logger.Printf("Error decoding person: %v", err)
			continue
		}
		people = append(people, &person)

	}
	// wg.Wait()
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, err := response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		if err != nil {
			logger.Printf("Error while writing the error response: %v\n", err)
		}
		return nil
	}

	return people

}
