// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// //`json: locationName`
// type weatherData struct {
// 	LocationName string   `json:"locationName,omitempty"`
// 	Weather      string   `json:"weather,omitempty"`
// 	Temperature  int      `json:"temperature,omitempty"`
// 	Celsius      bool     `json:"celsius.omitempty"`
// 	TempForecast []int    `json:"temp_forecast,omitempty"`
// 	Wind         windData `json:"wind,omitempty"`
// }

// type windData struct {
// 	Direction string `json:"direction,omitempty"`
// 	Speed     int    `json:"speed,omitempty"`
// }

// type loc struct {
// 	Lat float32 `json:"lat,omitempty"`
// 	Lon float32 `json:"lon,omitempty"`
// }

// func weatherHandler(w http.ResponseWriter, r *http.Request) {

// 	location := loc{}

// 	jsn, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		log.Fatal("Error reading the body", err)
// 	}

// 	err = json.Unmarshal(jsn, &location)
// 	if err != nil {
// 		log.Fatal("Decoding error: ", err)
// 	}

// 	log.Printf("Received: %v\n", location)

// 	weather := weatherData{
// 		LocationName: "Zzyzx",
// 		Weather:      "cloudy",
// 		Temperature:  31,
// 		Celsius:      true,
// 		TempForecast: []int{30, 32, 29},
// 		Wind: windData{
// 			Direction: "S",
// 			Speed:     20,
// 		},
// 	}

// 	weatherJson, err := json.Marshal(weather)
// 	if err != nil {
// 		fmt.Fprintf(w, "Error: %s", err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	w.Write(weatherJson)

// }

// func server() {
// 	http.HandleFunc("/weather", weatherHandler)
// 	http.ListenAndServe(":8088", nil)
// }

// func client() {

// 	locJson, err := json.Marshal(loc{Lat: 35.14326, Lon: -116.104})

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	req, err1 := http.NewRequest("POST", "http://localhost:8088/weather", bytes.NewBuffer(locJson))
// 	req.Header.Set("Content-Type", "application/json")

// 	if err1 != nil {
// 		log.Fatal("this was the error", err)
// 	}

// 	client := &http.Client{}
// 	resp, err2 := client.Do(req)
// 	if err2 != nil {
// 		log.Fatal("error in the client.do ", err)
// 	}
// 	body, err3 := ioutil.ReadAll(resp.Body)

// 	if err3 != nil {
// 		log.Printf("Error reading response body: %v", err1)
// 	}

// 	fmt.Println("Response: ", string(body))
// 	fmt.Println("status: ", resp.Status)
// 	resp.Body.Close()
// }

// func main() {

// 	go server()
// 	// fmt.Println("we are after server")
// 	client()
// }
