package main

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/connection"
	"go-rest-api/model"
	"net/http"
)

func main() {
	// Initialize the database connection
	connection.DbConnect()
	router := gin.Default()

	// Define a route that returns a JSON response
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	router.GET("/events", func(c *gin.Context) {
		getAllEvents := model.GetAllEvents()
		c.JSON(http.StatusOK, getAllEvents)
	})

	router.POST("/events", func(c *gin.Context) {
		var createEvent model.Event
		err := c.BindJSON(&createEvent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		createEvent.Id = 1
		createEvent.UserIds = 2
		createEvent.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Event created successfully!", "event": createEvent})
	})

	// Start the server on port 8080
	if err := router.Run(":3000"); err != nil {
		panic(err)
	}

}
