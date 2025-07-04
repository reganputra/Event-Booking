package main

import (
	"go-rest-api/connection"
	"go-rest-api/controllers"
	"go-rest-api/helper"
	"go-rest-api/middleware"
	"go-rest-api/repository"
	"go-rest-api/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default environment variables")
	}

	// Initialize the database connection
	db := connection.DbConnect()
	defer db.Close()

	// Initialize the repository
	eventRepo := repository.NewEventRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize the service
	eventService := services.NewEventService(eventRepo)
	userService := services.NewUserService(userRepo)

	// Initialize the controller
	eventController := controllers.NewEventController(eventService)
	userController := controllers.NewUserController(userService)

	router := gin.Default()

	// Healthcheck endpoint to verify server status
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is running!",
		})
	})

	// Public routes
	router.GET("/events", eventController.GetAllEvents)
	router.GET("/events/category/:category", eventController.GetEventsByCategory)
	router.GET("/events/:id", eventController.GetEventByID)
	router.POST("/users/register", userController.RegisterUser)
	router.POST("/users/login", userController.LoginUser)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.POST("/events", eventController.CreateEvent)
		protectedRoutes.PUT("/events/:id", eventController.UpdateEvent)
		protectedRoutes.DELETE("/events/:id", eventController.DeleteEvent)
		protectedRoutes.POST("/events/:id/register", eventController.RegisterForEvent)
		protectedRoutes.DELETE("/events/:id/register", eventController.CancelEventRegistration)
		protectedRoutes.GET("/events/registered", eventController.GetRegisteredEvents)
	}

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(middleware.AuthorizeRole("admin"))
	{
		adminRoutes.GET("/users", userController.GetAllUser)
		adminRoutes.GET("/users/:id", userController.GetUserByID)
		adminRoutes.PUT("/users/:id", userController.UpdateUser)
		adminRoutes.DELETE("/users/:id", userController.DeleteUser)
	}

	// Start the server on port 3000
	err = router.Run(":3000")
	helper.PanicIfError(err)

}
