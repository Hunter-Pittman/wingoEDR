package server

import (
	"io"
	"net/http"
	"os"
	"wingoEDR/usermanagement"
	"wingoEDR/webserver/models"

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

	superGroup := router.Group("/api")
	{
		superGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		managementGroup := superGroup.Group("/management")
		{
			wingoManagementGroup := managementGroup.Group("/wingo-management")
			{
				wingoManagementGroup.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "pong",
					})
				})
			}

			userManagementGroup := managementGroup.Group("/user-management")
			{
				userManagementGroup.GET("/users", func(c *gin.Context) {
					users := usermanagement.ReturnUsers()
					c.JSON(http.StatusOK, users)
				})

				userManagementGroup.POST("/create", func(c *gin.Context) {
					var user models.NewUser

					// bind the request body to the User struct
					if err := c.ShouldBindJSON(&user); err != nil {
						c.JSON(400, gin.H{"error": err.Error()})
						return
					}

					err := models.CreateNewUser(user)
					if err != nil {
						c.JSON(400, gin.H{"error": err.Error()})
						return
					}

					// return a success response
					c.JSON(200, gin.H{
						"message": "User created successfully",
						"user":    user,
					})
				})
			}

		}
	}

	return router
}
