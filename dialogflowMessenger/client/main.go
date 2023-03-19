package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"

	"src/github.com/DavidHernandez21/dialogflowMessenger/tokenGenerator"

	"github.com/google/uuid"
)

// projects/acn-aibot-stg/locations/europe-west1/agents/9ab6d4ba-8e5d-4567-b086-5415aaa6a2be

func main() {

	var jsonData = []byte(`{
		"queryInput": {
		  "text": {
			"text": "Ciao"
		  },
		  "languageCode": "it"
		},
		"queryParams": {
		  "timeZone": "Europe/Rome"
		}
	  }`)

	cmd, err := exec.Command("gcloud", "auth", "application-default", "print-access-token").Output()

	if err != nil {
		log.Fatalf("error executing the gcloud auth command: %v", err)
	}

	// log.Println(string(cmd))
	uuidWithHyphen := uuid.New().String()
	httpposturl := fmt.Sprintf("https://europe-west1-dialogflow.googleapis.com/v3/projects/acn-aibot-stg/locations/europe-west1/agents/9ab6d4ba-8e5d-4567-b086-5415aaa6a2be/sessions/%v:detectIntent", uuidWithHyphen)
	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("error creating a post request: %v", err)
	}

	var scope string = "https://www.googleapis.com/auth/dialogflow"

	token, err := tokenGenerator.GenerateTokenString(request.Context(), scope)

	if err != nil {
		log.Fatalf("error generating the token: %v", err)
	}

	re := regexp.MustCompile(`\r?\n`)
	re.ReplaceAllString(string(cmd), "")

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	request.Header.Set("X-Goog-User-Project", "acn-aibot-stg")
	// "X-Goog-User-Project"

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("post request error: %v", err)
	}
	defer response.Body.Close()

	log.Println("response Status:", response.Status)
	log.Println("response Headers:", response.Header)
	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("error reading the response: %v", err)
	}
	log.Println("response Body:", string(body))
}
