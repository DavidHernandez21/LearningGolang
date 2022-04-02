package main

import (
	"github.com/gin-gonic/gin"
)

func endpointHandler(c *gin.Context) {
	c.String(200, "%s %s", c.Request.Method, c.Request.URL.Path)
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"92.244.24.19"})
	router.GET("/products", endpointHandler)
	router.GET("/products/:productId", endpointHandler)
	// Eg: /products/1052
	router.POST("/products", endpointHandler)
	router.PUT("/products/:productId", endpointHandler)
	router.DELETE("/products/:productId", endpointHandler)
	router.Run("127.0.0.1:5000")
}
