package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/DavidHernandez21/gingonic/getserver"
	"github.com/DavidHernandez21/gingonic/servershutdown"
	"github.com/gin-gonic/gin"
)

type AddParams struct {
	X float64 `json:"x" binding:"required,number,gt=10"`
	Y float64 `json:"y" binding:"required,numeric"`
}

func add(c *gin.Context) {
	var ap AddParams
	if err := c.ShouldBindJSON(&ap); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"answer": ap.X + ap.Y})
}

func main() {
	bindAddress := flag.String("bind", "127.0.0.1:8080", "Bind address")
	flag.Parse()

	router := gin.Default()
	router.SetTrustedProxies([]string{"92.244.24.19"})
	router.POST("/add", add)

	srv := getserver.NewServer(*bindAddress, router)

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	servershutdown.Graceful(srv, 5*time.Second)

}
