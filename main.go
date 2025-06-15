package main

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/connection"
	"go-rest-api/controllers"
	"go-rest-api/helper"
	"go-rest-api/middleware"
	"go-rest-api/repository"
	"go-rest-api/services"
	"net/http"
)

func main() {
	// Initialize the database connection
	db := connection.DbConnect()
	defer db.Close()

	// Initialize the event repository
	eventRepo := repository.NewEventRepository(db)

	// Intialize the event service
	eventService := services.NewEventService(eventRepo)

	// Initialize the event controller
	eventController := controllers.NewEventController(eventService)

	router := gin.Default()

	// Healthcheck endpoint to verify server status
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is running!",
		})
	})

	// Public routes
	router.GET("/events", eventController.GetAllEvents)
	router.GET("/events/:id", eventController.GetEventByID)
	router.POST("/users/register", controllers.RegisterUser)
	router.POST("/users/login", controllers.LoginUser)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.POST("/events", eventController.CreateEvent)
		protectedRoutes.PUT("/events/:id", eventController.UpdateEvent)
		protectedRoutes.DELETE("/events/:id", eventController.DeleteEvent)
		protectedRoutes.POST("/events/:id/register", eventController.RegisterForEvent)
		protectedRoutes.DELETE("/events/:id/register", eventController.CancelEventRegistration)
	}

	// Start the server on port 3000
	err := router.Run(":3000")
	helper.PanicIfError(err)

}
