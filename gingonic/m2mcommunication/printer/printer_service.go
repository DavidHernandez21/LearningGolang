package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/DavidHernandez21/gingonic/getserver"
	"github.com/DavidHernandez21/gingonic/servershutdown"
	"github.com/gin-gonic/gin"
)

type PrintJob struct {
	Format    string `json:"format" binding:"required"`
	InvoiceId int    `json:"invoiceId" binding:"required,gte=0"`
	JobId     int    `json:"jobId" binding:"gte=0"`
}

func main() {
	bindAddress := flag.String("bind", "127.0.0.1:8080", "Bind address")
	flag.Parse()

	router := gin.Default()
	router.SetTrustedProxies([]string{"92.244.24.19"})

	router.POST("/print-jobs", func(c *gin.Context) {
		var p PrintJob
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input!"})
			return
		}
		log.Printf("PrintService: creating new print job from invoice #%v...", p.InvoiceId)
		rand.Seed(time.Now().UnixNano())
		p.JobId = rand.Intn(1000)
		log.Printf("PrintService: created print job #%v", p.JobId)
		c.JSON(200, p)
	})

	srv := getserver.NewServer(*bindAddress, router)

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	servershutdown.Graceful(&srv, 5*time.Second)
}
