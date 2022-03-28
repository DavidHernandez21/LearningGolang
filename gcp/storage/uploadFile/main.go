package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// type Token struct {
// 	// AccessToken is the token that authorizes and authenticates
// 	// the requests.
// 	AccessToken string `json:"access_token"`

// 	// TokenType is the type of token.
// 	// The Type method returns either this or "Bearer", the default.
// 	TokenType string `json:"token_type,omitempty"`

// 	// RefreshToken is a token that's used by the application
// 	// (as opposed to the user) to refresh the access token
// 	// if it expires.
// 	RefreshToken string `json:"refresh_token,omitempty"`

// 	// Expiry is the optional expiration time of the access token.
// 	//
// 	// If zero, TokenSource implementations will reuse the same
// 	// token forever and RefreshToken or equivalent
// 	// mechanisms for that TokenSource will not be used.
// 	Expiry time.Time `json:"expiry,omitempty"`
// 	// contains filtered or unexported fields
// }

// func (t *Token) Token() (*Token, error) {

// 	cmd, err := exec.Command("gcloud", "auth", "application-default", "print-access-token").Output()

// 	if err != nil {
// 		return nil, fmt.Errorf("error executing the gcloud auth command: %v", err)
// 	}

// 	re := regexp.MustCompile(`\r?\n`)
// 	tokenString := re.ReplaceAllString(string(cmd), "")

// 	t.AccessToken = tokenString
// 	t.Expiry = time.Now().Add(60 * time.Minute)

// 	return t, nil

// }

func uploadFileToBucket(bucket, object, fileName, jsonkeyPath string) error {

	ctx := context.Background()

	// var token *Token

	// token, err := token.Token()

	// if err != nil {
	// 	return fmt.Errorf("error generating token: %v", err)
	// }
	// inputJsonKeyPath := flag.String("", "key.json", "json key file path")
	// outputDirectory := flag.String("0", "./", "output directory")

	// flag.Parse()

	// b, err := ioutil.ReadFile(jsonkeyPath)

	// if err != nil {
	// 	return fmt.Errorf("error reading json key: %v", err)
	// }
	// tokenGenerator.Keyfile = jsonkeyPath
	// // fmt.Printf("keyfile: %v\n", tokenGenerator.Keyfile)
	// var scope string = "https://www.googleapis.com/auth/devstorage.read_write"
	// token, err := tokenGenerator.GenerateToken(context.Background(), scope)

	// if err != nil {
	// 	return fmt.Errorf("error generating token: %v", err)
	// }

	// option.WithTokenSource(token)
	// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", jsonkeyPath)
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonkeyPath))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)

	wc.ObjectAttrs.ContentType = "text/csv"
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(os.Stdout, "Blob %v uploaded.\n", object)

	// if err := os.Remove(f.Name()); err != nil {
	// 	return fmt.Errorf("deleting file '%v': %v", f.Name(), err)
	// }
	return nil
}

func main() {
	inputFilePath := flag.String("i", "DATA_GIA_GOOGLE_20201023-161045_prod.csv", "input file path")
	// outputDirectory := flag.String("o", "./", "output directory")
	bucketName := flag.String("bucket", "gpjegld01-1812-convintel-ivr-csv-tr-gcf", "name of the gcs bucket")
	inputJsonKeyPath := flag.String("jsonKey", "key.json", "json key file path")

	flag.Parse()

	err := uploadFileToBucket(*bucketName, "test.csv", *inputFilePath, *inputJsonKeyPath)

	if err != nil {
		fmt.Printf("error uploading file '%v' to gcs bucket '%v': %v", *inputFilePath, *bucketName, err)
	}

}
