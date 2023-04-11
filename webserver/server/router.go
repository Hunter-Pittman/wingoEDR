package server

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	// Disable Console Color, you don't need console color when writing the logs to file.
	//gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	//router.Use(auth.TokenAuthMiddleware(apiToken))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}
