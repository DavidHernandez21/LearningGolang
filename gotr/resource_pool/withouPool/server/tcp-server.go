package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"src/github.com/DavidHernandez21/gotr/resource_pool/withouPool/model"
	"sync/atomic"
	"time"

	log "github.com/mgutz/logxi/v1"
)

type (
	// TCPServer listen for client on a port and accepts messages
	TCPServer struct {
		numReqs uint64
		port    string
		log     log.Logger
		s       *http.Server
	}
)

var (
	s  = rand.NewSource(time.Now().Unix())
	rn = rand.New(s)
)

// newTCPServer creates a new server to accept
func newTCPServer(port string) Server {
	srv := &TCPServer{port: port}
	srv.log = log.New("server")
	srv.log.SetLevel(log.LevelInfo)

	// TODO - create and configure http.Server
	s := &http.Server{}
	// configure http server
	s.WriteTimeout = 500 * time.Millisecond
	s.ReadTimeout = 1000 * time.Millisecond
	s.Addr = fmt.Sprintf("localhost:%v", srv.port) // listen on all interfaces
	s.Handler = srv
	srv.s = s // store a reference to the http.Server
	return srv
}

// Start listening on srv.port and all interfaces
func (srv *TCPServer) Start() error {
	if nil == srv {
		return fmt.Errorf("Start() called on nil TCPServer object")
	}

	srv.log.Info("Starting HTTP server, listening on port", srv.port)
	// TODO - actually start the http server
	var err error
	go func() {
		err = srv.s.ListenAndServe()
	}()
	time.Sleep(200 * time.Millisecond)
	return err
}

// Stop listening and close all client connections
func (srv *TCPServer) Stop() {
	if nil == srv {
		return
	}

	srv.log.Info("Stopping HTTP server")
	// TODO - stop http server
	defer srv.s.Close()

	srv.log.Info("Messages processed:", srv.numReqs)
}

func (srv *TCPServer) Shutdown(ctx context.Context) error {
	if nil == srv {
		return nil
	}

	srv.log.Info("Trying to shutdown gracefully")
	if err := srv.s.Shutdown(ctx); err != nil {
		return err
	}

	srv.log.Info("Messages processed:", srv.numReqs)

	return nil

}

func (srv *TCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&srv.numReqs, 1)

	srv.log.Debug("TCPServer - message from", r.RemoteAddr)
	// this is already called from a goroutine for each connection
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	msg := &model.ClientReq{}
	err := dec.Decode(msg)
	if err != nil {
		srv.log.Debug("Error decoding message", err)
	}
	// INFO - pretent we do some work on with the msg
	time.Sleep(time.Duration(rn.Int63n(5)) * time.Millisecond)
}
