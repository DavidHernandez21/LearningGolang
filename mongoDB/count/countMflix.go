package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	// uri = os.Getenv("MONGODB_URI")
	uri = os.Getenv("MONGODB_URI_WO_DATABASE")
	uri = strings.Replace(uri, "<Database>", "sample_mflix", 1)

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// begin countDocuments
	// coll := client.Database("sample_mflix").Collection("movies")
	coll := client.Database("thepolyglotdeveloper").Collection("people")
	filter := bson.D{primitive.E{Key: "countries", Value: "China"}}

	ctxFullCount, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	estCount, estCountErr := coll.EstimatedDocumentCount(ctxFullCount)
	count, countErr := coll.CountDocuments(ctxFullCount, filter)
	// end countDocuments

	if estCountErr != nil {
		panic(estCountErr)
	}
	if countErr != nil {
		panic(countErr)
	}

	// When you run this file, it should print:
	// Estimated number of documents in the movies collection: 23541
	// Number of movies from China: 303
	fmt.Printf("Estimated number of documents in the movies collection: %d\n", estCount)
	fmt.Printf("Number of movies from China: %d\n", count)
}
