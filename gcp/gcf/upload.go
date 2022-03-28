package p

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// upload file to a bucket, the file to upload comes from a channel. After upload the file sends the name of the file to a channel
func uploadFileToBucket(bucket, fileName string, client *storage.Client) error {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	object := fileName[2:]

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("filename: %v, os.Open: %v", fileName, err)
	}

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)

	wc.ObjectAttrs.ContentType = "text/csv"
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("filename: %v, io.Copy: %v", object, err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("filename: %v, Writer.Close: %v", object, err)
	}
	fmt.Fprintf(os.Stdout, "Blob %v uploaded to bucket %v\n", object, bucket)

	if err := f.Close(); err != nil {
		return fmt.Errorf("filename: %v, File.Close: %v", object, err)

	}

	return nil

}
