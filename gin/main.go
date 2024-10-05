package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	ginServer := gin.Default()
	ginServer.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	ginServer.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
