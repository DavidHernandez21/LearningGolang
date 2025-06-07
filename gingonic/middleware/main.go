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

// look for resourcer on this link https://github.com/gin-contrib

func FindUserAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.GetHeader("User-Agent"))
		// Before calling handler
		c.Next()
		// After calling handler
	}
}
func main() {
	bindAddress := flag.String("bind", "127.0.0.1:8080", "Bind address")
	flag.Parse()

	router := gin.Default()
	router.SetTrustedProxies([]string{"92.244.24.19"})

	router.Use(FindUserAgent())
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Middleware works!"})
	})

	srv := getserver.NewServer(*bindAddress, router)

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	servershutdown.Graceful(&srv, 5*time.Second)

}
