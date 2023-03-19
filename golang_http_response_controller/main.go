package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter // Embed a http.ResponseWriter
	statusCode          int
	headerWritten       bool
}

func newstatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (mw *statusResponseWriter) WriteHeader(statusCode int) {
	mw.ResponseWriter.WriteHeader(statusCode)

	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

func (mw *statusResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true
	return mw.ResponseWriter.Write(b)
}

func (mw *statusResponseWriter) Unwrap() http.ResponseWriter {
	return mw.ResponseWriter
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/normal", normalHandler)
	mux.HandleFunc("/flushed", flushedHandler)

	const port = 3000

	log.Printf("Listening on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), logResponse(mux))
	if err != nil {
		log.Fatal(err)
	}
}

func logResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := newstatusResponseWriter(w)
		next.ServeHTTP(sw, r)
		log.Printf("%s %s: status %d\n", r.Method, r.URL.Path, sw.statusCode)
	})
}

func normalHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("OK"))
}

func flushedHandler(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)

	w.Write([]byte("Write A...."))
	err := rc.Flush()
	if err != nil {
		log.Println(err)
		return
	}

	time.Sleep(time.Second)

	w.Write([]byte("Write B...."))
	err = rc.Flush()
	if err != nil {
		log.Println(err)
	}
}
