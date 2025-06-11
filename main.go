package main

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/connection"
	"go-rest-api/controllers"
	"go-rest-api/helper"
	"net/http"
)

func main() {
	// Initialize the database connection
	connection.DbConnect()
	router := gin.Default()

	// Healthcheck endpoint to verify server status
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is running!",
		})
	})

	// Router for event-related operations
	router.GET("/events", controllers.GetAllEvents)
	router.GET("/events/:id", controllers.GetEventsById)
	router.POST("/events", controllers.CreateEvent)
	router.PUT("/events/:id", controllers.UpdateEvent)
	router.DELETE("/events/:id", controllers.DeleteEvent)

	// Router for user-related operations
	router.POST("/users/register", controllers.RegisterUser)

	// Start the server on port 3000
	err := router.Run(":3000")
	helper.PanicIfError(err)

}
