package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	rg = rand.New(rand.NewSource(time.Now().Unix()))
)

func main() {
	client := http.DefaultClient
	var errBody error

	for i := 0; i < 10; {
		resp, err := client.Get("http://localhost:8080/api")
		if err == nil && resp.StatusCode == http.StatusOK {
			i++
			io.Copy(os.Stdout, resp.Body)
			errBody = resp.Body.Close()
			if errBody != nil {
				logrus.Errorf("Error closing body: %v", errBody)
			}
			continue
		}
		logrus.Errorf("request failed, status code: %v", resp.StatusCode)
		d := time.Duration(rg.Int31n(500)) * time.Millisecond
		errBody = resp.Body.Close()
		if err != nil {
			logrus.Errorf("Error closing body: %v", errBody)
		}
		time.Sleep(d)

	}
}
