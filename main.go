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
	reviewRepo := repository.NewReviewRepository(db)     // Add ReviewRepository
	waitlistRepo := repository.NewWaitlistRepository(db) // Add WaitlistRepository

	// Initialize the service
	// Note: WaitlistService might need UserRepository if it fetches user details for notifications
	waitlistService := services.NewWaitlistService(waitlistRepo, eventRepo, userRepo)
	eventService := services.NewEventService(eventRepo, waitlistService) // Pass waitlistService to EventService
	userService := services.NewUserService(userRepo)
	reviewService := services.NewReviewService(reviewRepo, eventRepo)

	// Initialize the controller
	eventController := controllers.NewEventController(eventService)
	userController := controllers.NewUserController(userService)
	reviewController := controllers.NewReviewController(reviewService)
	waitlistController := controllers.NewWaitlistController(waitlistService) // Add WaitlistController

	router := gin.Default()

	// Healthcheck endpoint to verify server status
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is running!",
		})
	})

	// Public routes
	router.GET("/events", eventController.GetAllEvents)
	router.GET("/events/search", eventController.SearchEvents) // New search endpoint
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

		protectedRoutes.POST("/events/:id/reviews", reviewController.CreateReview)

		// Waitlist routes (Protected)
		protectedRoutes.POST("/events/:id/waitlist", waitlistController.JoinWaitlist)
		protectedRoutes.DELETE("/events/:id/waitlist", waitlistController.LeaveWaitlist)
		// GET /events/:id/waitlist is admin/owner only - requires more specific authorization
		// For now, placing it under general protected routes. Access control needs to be refined in the controller or via specific middleware.
		// A simple approach for now: only admins can view waitlists.
		// If event owners should also view, the AuthorizeRole("admin") middleware needs to be more flexible or applied differently.
	}
	// Public route for getting reviews for an event
	router.GET("/events/:id/reviews", reviewController.GetReviewsForEvent)

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(middleware.AuthorizeRole("admin"))
	{
		adminRoutes.GET("/users", userController.GetAllUser)
		adminRoutes.GET("/users/:id", userController.GetUserByID)
		adminRoutes.PUT("/users/:id", userController.UpdateUser)
		adminRoutes.DELETE("/users/:id", userController.DeleteUser)

		// Waitlist viewing for admin
		adminRoutes.GET("/events/:id/waitlist", waitlistController.GetWaitlistForEvent)
	}

	// Start the server on port 3000
	err = router.Run(":3000")
	helper.PanicIfError(err)

}
