package tokenGenerator

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type JsonKey struct {
	Type                    string `json:"type,omitempty"`
	ProjectID               string `json:"project_id,omitempty"`
	PrivateKeyID            string `json:"private_key_id,omitempty"`
	PrivateKey              string `json:"private_key,omitempty"`
	ClientEmail             string `json:"client_email,omitempty"`
	ClientID                string `json:"client_id,omitempty"`
	AuthUri                 string `json:"auth_uri,omitempty"`
	TokenUri                string `json:"token_uri,omitempty"`
	AuthProviderx509CertUrl string `json:"auth_provider_x509_cert_url,omitempty"`
	Clientx509CertUrl       string `json:"client_x509_cert_url,omitempty"`
}

var (
	// projectID = flag.String("project", "acn-aibot-stg", "Project ID")

	Keyfile = flag.String("keyfile", "C:/Users/asbel.hernandez/OneDrive - Accenture/Personal/vs_studio/golang/src/github.com/DavidHernandez21/dialogflowMessenger/tokenGenerator/acn-aibot-stg-3d82eb27842b.json", "Service Account JSON keyfile")
)

func GenerateTokenString(ctx context.Context, scope string) (string, error) {
	flag.Parse()

	// if *projectID == "" {
	// 	fmt.Fprintln(os.Stderr, "missing -project flag")
	// 	flag.Usage()
	// 	os.Exit(2)
	// }
	if *Keyfile == "" {
		fmt.Fprintln(os.Stderr, "missing -keyfile flag")
		flag.Usage()
		os.Exit(2)
	}

	var jsonKey JsonKey

	// audience values for other services can be found in the repo here similar to
	// PubSub
	// https://github.com/googleapis/googleapis/blob/master/google/pubsub/pubsub.yaml
	// var scope string = "https://www.googleapis.com/auth/dialogflow"

	// ctx := context.Background()
	keyBytes, err := os.ReadFile(*Keyfile)
	if err != nil {
		return "", fmt.Errorf("unable to read service account key file  %v", err)
	}

	if err := json.Unmarshal(keyBytes, &jsonKey); err != nil {
		return "", fmt.Errorf("error unmarshalling the file %v: %v", *Keyfile, err)
	}

	// start := time.Now()
	tokenSource, err := google.JWTConfigFromJSON(keyBytes, scope)
	if err != nil {
		return "", fmt.Errorf("error building JWT access token source: %v", err)
	}

	tokenSource.Audience = "https://www.googleapis.com/oauth2/v4/token"
	// log.Println(tokenSource.Audience)
	tokenSource.PrivateKeyID = jsonKey.PrivateKeyID
	tokenSource.TokenURL = jsonKey.TokenUri
	tokenSource.Email = jsonKey.ClientEmail
	tokenSource.PrivateKey = []byte(jsonKey.PrivateKey)
	// for _, v := range tokenSource.Scopes {
	// 	log.Println(v)
	// }
	// log.Printf("%v", tokenSource.UseIDToken)
	tokenSource.Expires = 3600 * time.Second
	// log.Printf("%v", tokenSource.Expires)

	jwt, err := tokenSource.TokenSource(ctx).Token()
	if err != nil {
		return "", fmt.Errorf("unable to generate JWT token: %v", err)
	}
	// log.Println(jwt.AccessToken)
	return jwt.AccessToken, nil

}

func GenerateToken(ctx context.Context, scope string) (oauth2.TokenSource, error) {
	flag.Parse()

	// if *projectID == "" {
	// 	fmt.Fprintln(os.Stderr, "missing -project flag")
	// 	flag.Usage()
	// 	os.Exit(2)
	// }
	if *Keyfile == "" {
		fmt.Fprintln(os.Stderr, "missing -keyfile flag")
		flag.Usage()
		os.Exit(2)
	}

	// fmt.Fprintf(os.Stdout, "keyfile: %v", *Keyfile)
	var jsonKey JsonKey

	// audience values for other services can be found in the repo here similar to
	// PubSub
	// https://github.com/googleapis/googleapis/blob/master/google/pubsub/pubsub.yaml
	// var scope string = "https://www.googleapis.com/auth/dialogflow"

	// ctx := context.Background()
	keyBytes, err := os.ReadFile(*Keyfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read service account key file  %v", err)
	}

	if err := json.Unmarshal(keyBytes, &jsonKey); err != nil {
		return nil, fmt.Errorf("error unmarshalling the file %v: %v", *Keyfile, err)
	}

	// start := time.Now()
	tokenSource, err := google.JWTConfigFromJSON(keyBytes, scope)
	if err != nil {
		return nil, fmt.Errorf("error building JWT access token source: %v", err)
	}

	tokenSource.Audience = "https://www.googleapis.com/oauth2/v4/token"
	// log.Println(tokenSource.Audience)
	tokenSource.PrivateKeyID = jsonKey.PrivateKeyID
	tokenSource.TokenURL = jsonKey.TokenUri
	tokenSource.Email = jsonKey.ClientEmail
	tokenSource.PrivateKey = []byte(jsonKey.PrivateKey)
	// for _, v := range tokenSource.Scopes {
	// 	log.Println(v)
	// }
	// log.Printf("%v", tokenSource.UseIDToken)
	tokenSource.Expires = 3600 * time.Second
	// log.Printf("%v", tokenSource.Expires)

	jwt := tokenSource.TokenSource(ctx)

	return jwt, nil
	// log.Println(jwt.AccessToken)
	// return jwt.AccessToken, nil

}
