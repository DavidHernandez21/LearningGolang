// package main

// import (
// 	"fmt"
// 	"os"
// 	"strings"
// )

// func main() {
// 	// Set custom env variable
// 	//   os.Setenv("CUSTOM", "500")

// 	// fetcha all env variables
// 	for _, element := range os.Environ() {
// 		variable := strings.Split(element, "=")
// 		if variable[0] == "MONGODB_URI" {
// 			fmt.Println("Found it ----------------")
// 			fmt.Println(variable[0], "=>", variable[1])
// 		}
// 		if variable[0] == "SPARK_HOME" {
// 			fmt.Println("Found it ----------------")
// 			fmt.Println(variable[0], "=>", variable[1])
// 		}
// 		// fmt.Println(variable[0], "=>", variable[1])
// 	}
// 	// s := "mongodb+srv://David7:uLAOvPHtGbR50aud@clusterdavid.mk1fu.mongodb.net/<Database>?retryWrites=true&w=majority"

// 	// s1 := strings.Replace(s, "<Database>", "Daje", 1)

// 	// fmt.Println(s1)
// }

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

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// uri := os.Getenv("MONGODB_URI")
	uri := os.Getenv("MONGODB_URI_WO_DATABASE")
	uri = strings.Replace(uri, "<Database>", "sample_mflix", 1)
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		cancel()
	}()
	coll := client.Database("sample_mflix").Collection("movies")
	title := "Back to the Future"
	var result bson.M
	err = coll.FindOne(ctx, bson.D{primitive.E{Key: "title", Value: title}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}
