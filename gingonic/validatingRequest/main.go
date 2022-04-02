package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DavidHernandez21/gingonic/getserver"
	servershutdown "github.com/DavidHernandez21/gingonic/serverShutdown"
	"github.com/gin-gonic/gin"
)

type PrintJob struct {
	JobId int `json:"jobId" binding:"required,gte=10000"`
	Pages int `json:"pages" binding:"required,gte=1,lte=100"`
}

func main() {
	bindAddress := flag.String("bind", "127.0.0.1:8080", "Bind address")
	flag.Parse()

	router := gin.Default()
	router.SetTrustedProxies([]string{"92.244.24.19"})
	router.POST("/print", func(c *gin.Context) {
		var p PrintJob
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input!"})
			return
		}
		c.JSON(200, gin.H{"message": fmt.Sprintf("PrintJob #%v started!", p.JobId)})
	})

	srv := getserver.NewServer(*bindAddress, router)

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	servershutdown.Graceful(srv, 5*time.Second)

}
