package hit

import (
	"io"
	"log"
	"net/http"
	"time"
)

// TIP Why export the Send function when only the do is using it?
// You do it to let the library's importers build
// their custom clients (as you did by creating the Client type).

// SendFunc is the type of the function called by Client.Do
// to send an HTTP request and return a performance result.
type SendFunc func(*http.Request) *Result

// var once sync.Once

// Send an HTTP request and return a performance result.
func Send(client *http.Client, r *http.Request) *Result {
	t := time.Now()

	var (
		code  uint16
		bytes int64
	)
	response, err := client.Do(r)

	if err == nil {
		code = uint16(response.StatusCode)
		bytes, err = io.Copy(io.Discard, response.Body)
		errClose := response.Body.Close()
		if errClose != nil {
			log.Printf("request to %s, error closing response body: %v", r.URL, errClose)
		}
	}

	return &Result{
		Duration: time.Since(t),
		Bytes:    bytes,
		Status:   code,
		Error:    err,
	}
}

// func SendSlow(r *http.Request) *Result {
// 	t := time.Now()

// 	var (
// 		code  uint16
// 		bytes int64
// 	)
// 	response, err := http.DefaultClient.Do(r)

// 	if err == nil {
// 		code = uint16(response.StatusCode)
// 		bytes, err = io.Copy(io.Discard, response.Body)
// 		errClose := response.Body.Close()
// 		if errClose != nil {
// 			log.Printf("request to %s, error closing response body: %v", r.URL, errClose)
// 		}
// 	}

// 	return &Result{
// 		Duration: time.Since(t),
// 		Bytes:    bytes,
// 		Status:   code,
// 		Error:    err,
// 	}
// }
