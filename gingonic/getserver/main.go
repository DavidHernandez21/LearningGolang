package getserver

import (
	"net/http"
	"time"
)

// Returns a new Server instance with the supplied handler and address.
func NewServer(address string, handler http.Handler) http.Server {
	return http.Server{
		Addr:    address, // configure the bind address
		Handler: handler, // set the default handler
		// ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
}
