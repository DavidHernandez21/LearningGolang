// Package p contains an HTTP Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.

var bucket = os.Getenv("BUCKET_NAME")

// client is a global storage client, initialized once per instance.
var client *storage.Client

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	// client is initialized with context.Background() because it should
	// persist between function invocations.
	client, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("error creating storage.NewClient: %v", err)
	}
}

func SplitCsv(inputFilePath, outputDirectory, bucket string, client *storage.Client) error {
	splitter := NewSplitter()
	splitter.FileChunkSize = 500000 //in bytes (50MB)
	_, err := splitter.Split(inputFilePath, outputDirectory, bucket, client)
	if err != nil {
		return fmt.Errorf("error in the SplitCsv function: %v", err)
	}
	return nil

}

func changeDir(directory string, w http.ResponseWriter) {
	log.Printf("changing directory to %v", directory)
	if err := os.Chdir(directory); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("error while changing to directory '%v': %v\n", directory, err)
		return
	}
}

func SplitUploadAndDeleteFiles(w http.ResponseWriter, r *http.Request) {

	// log.Printf("request: %#v\n", &r.Body)

	var d struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {

		log.Printf("json.NewDecoder: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return

	}

	changeDir("/tmp", w)

	filename := d.Message
	if strings.TrimSpace(filename) == "" {
		log.Printf("Got an empty csv filename")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// "DATA_GIA_GOOGLE_20220103-080110_6.csv"
	err := downloadFile(bucket, filename, filename, client)

	if err != nil {
		log.Printf("error downloading file: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = SplitCsv(filename, "./", bucket, client)

	if err != nil {
		log.Printf("error splitting the files: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		deleteFiles(filename)
		return
	}

	deleteFiles(filename)
	changeDir("..", w)

	// log.Printf("struct %v", d)
	// if d.Message == "" {
	// 	fmt.Fprint(w, "Hello World!")
	// 	return
	// }
	fmt.Fprintf(w, "file %v correctly proccessed", filename)
}
