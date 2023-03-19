package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"src/github.com/DavidHernandez21/gotr/resource_pool/syncPool/model"
	"time"

	log "github.com/mgutz/logxi/v1"
)

const (
	requestsPerClient = 100000
	maxBatchSize      = (requestsPerClient / 10) * 2 // 20% of total request
)

var (
	s      = rand.NewSource(time.Now().Unix())
	r      = rand.New(s)
	logger = log.New("client")
	// buf    = &bytes.Buffer{}
)

func submitRequests(url string) {
	var req *model.ClientReq
	msgLeft := requestsPerClient
	var reqID uint

	for 0 < msgLeft {
		batch := r.Intn(maxBatchSize)
		if batch > msgLeft {
			batch = msgLeft
		}
		msgLeft -= batch

		for i := 0; i < batch; i++ {
			req = &model.ClientReq{}
			reqID++
			req.ID = reqID
			req.Size = r.Intn(model.ReqDataSize)
			buf, err := encodeReq(req)
			if nil != err {
				logger.Debug("JSON encode error:", err)
				break // try again later
			}
			// fmt.Println(buf) // send to server
			resp, err := http.Post(url, "text/json", buf)
			if nil != err {
				logger.Debug("Post error:", err)
				break // try again later
			}
			defer resp.Body.Close()
		}
		// pause a bit between batches
		// time.Sleep(time.Duration(r.Intn(200)) * time.Millisecond)
	}
}

func encodeReq(req *model.ClientReq) (io.Reader, error) {
	// buf.Reset()
	var buf = &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func StartTimer(name string) func() {
	t := time.Now()
	fmt.Println(name, "started")

	return func() {
		// d := time.Now().Sub(t)
		d := time.Since(t)
		fmt.Println(name, "took", d)
	}

}
