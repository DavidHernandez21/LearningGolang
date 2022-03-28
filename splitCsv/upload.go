package splitCsv

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func getClient(jsonkeyPath string) (*storage.Client, error) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonkeyPath))
	if err != nil {
		return nil, fmt.Errorf("could not create sotrage client: %v", err)
	}

	return client, nil

}

func uploadFileToBucket(bucket string, chIn <-chan string, chOut chan<- string, errChan chan error, wg *sync.WaitGroup, client *storage.Client) {

	go func() {
		defer wg.Done()
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		for fileName := range chIn {

			object := fileName[2:]

			f, err := os.Open(fileName)
			if err != nil {
				errChan <- fmt.Errorf("filename: %v, os.Open: %v", fileName, err)
			}

			// Upload an object with storage.Writer.
			wc := client.Bucket(bucket).Object(object).NewWriter(ctx)

			wc.ObjectAttrs.ContentType = "text/csv"
			if _, err = io.Copy(wc, f); err != nil {
				errChan <- fmt.Errorf("filename: %v, io.Copy: %v", object, err)
			}
			if err := wc.Close(); err != nil {
				errChan <- fmt.Errorf("filename: %v, Writer.Close: %v", object, err)
			}
			fmt.Fprintf(os.Stdout, "Blob %v uploaded to bucket %v\n", object, bucket)

			if err := f.Close(); err != nil {
				errChan <- fmt.Errorf("filename: %v, File.Close: %v", object, err)

				fmt.Fprintf(os.Stderr, "File NOT sent to delete, error closing file '%v': %v\n", object, err)
				continue
			}

			fmt.Fprintf(os.Stdout, "file send to delete: %v\n", fileName)
			chOut <- fileName

		}
	}()

	// close(chOut)
	// close(errChan)

}
