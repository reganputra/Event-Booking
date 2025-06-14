package main

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/connection"
	"go-rest-api/controllers"
	"go-rest-api/helper"
	"go-rest-api/middleware"
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

	// Public routes
	router.GET("/events", controllers.GetAllEvents)
	router.GET("/events/:id", controllers.GetEventsById)
	router.POST("/users/register", controllers.RegisterUser)
	router.POST("/users/login", controllers.LoginUser)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.POST("/events", controllers.CreateEvent)
		protectedRoutes.PUT("/events/:id", controllers.UpdateEvent)
		protectedRoutes.DELETE("/events/:id", controllers.DeleteEvent)
		protectedRoutes.POST("/events/:id/register", controllers.RegisterForEvent)
		protectedRoutes.DELETE("/events/:id/register", controllers.CancelForEvent)
	}

	// Start the server on port 3000
	err := router.Run(":3000")
	helper.PanicIfError(err)

}
