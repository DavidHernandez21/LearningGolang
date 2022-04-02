package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DavidHernandez21/gingonic/getserver"
	"github.com/DavidHernandez21/gingonic/servershutdown"
	"github.com/gin-gonic/gin"
)

func add(c *gin.Context) {
	x, err := strconv.ParseFloat(c.Param("x"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	y, err := strconv.ParseFloat(c.Param("y"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.String(200, fmt.Sprintf("%.2f", x+y))
}

func main() {
	bindAddress := flag.String("bind", "127.0.0.1:8080", "Bind address")
	flag.Parse()

	router := gin.Default()
	router.SetTrustedProxies([]string{"92.244.24.19"})
	router.GET("/add/:x/:y", add)

	// srv := http.Server{
	// 	Addr:    *bindAddress, // configure the bind address
	// 	Handler: router,       // set the default handler
	// 	// ErrorLog:     l,                 // set the logger for the server
	// 	ReadTimeout:  5 * time.Second,   // max time to read request from the client
	// 	WriteTimeout: 10 * time.Second,  // max time to write response to the client
	// 	IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	// }
	srv := getserver.NewServer(*bindAddress, router)

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	servershutdown.Graceful(srv, 5*time.Second)
}
