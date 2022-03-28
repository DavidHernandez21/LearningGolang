package main

import (
	"context"
	"encoding/json"
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

// begin struct
type BlogPost struct {
	Title       string
	Author      string
	WordCount   int `bson:"word_count,omitempty"`
	LastUpdated time.Time
	Tags        []string
}

// end struct

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	uri = os.Getenv("MONGODB_URI")
	uri = strings.Replace(uri, "<Database>", "sample_training", 1)

	if uri == "" {
		log.Fatal("You must set your `MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		cancel()
	}()

	// begin create and insert
	myCollection := client.Database("sample_training").Collection("posts")

	post := BlogPost{
		Title:       "Annuals vs. Perennials?",
		Author:      "Sam Lee",
		WordCount:   682,
		LastUpdated: time.Now(),
		Tags:        []string{"seasons", "gardening", "flower"},
	}

	_, err = myCollection.InsertOne(ctx, post)
	// end create and insert

	if err != nil {
		panic(err)
	}

	var result bson.M
	err = myCollection.FindOne(ctx, bson.D{primitive.E{Key: "author", Value: "Sam Lee"}}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Could not find any document authored by Sam Lee\n")
			return
		}
		panic(err)
	}

	output, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", output)
}
