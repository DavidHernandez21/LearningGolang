package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	uri = os.Getenv("MONGODB_URI")
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

	// begin distinct
	coll := client.Database("sample_mflix").Collection("movies")
	filter := bson.D{primitive.E{Key: "directors", Value: "Natalie Portman"}}

	results, err := coll.Distinct(ctx, "title", filter)
	// end distinct

	if err != nil {
		panic(err)
	}

	// When you run this file, it should print:
	// A Tale of Love and Darkness
	// New York, I Love You
	for _, result := range results {
		fmt.Println(result)
	}
}
