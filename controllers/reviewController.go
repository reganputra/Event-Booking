package controllers

import (
	"database/sql"
	"errors"
	"go-rest-api/model"
	"go-rest-api/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	reviewService services.ReviewService
}

func NewReviewController(reviewService services.ReviewService) *ReviewController {
	return &ReviewController{reviewService: reviewService}
}

func (c *ReviewController) CreateReview(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDVal.(int64)

	eventIDStr := ctx.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	var review model.Review
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.EventID = eventID // Set EventID from path parameter

	if err := c.reviewService.CreateReview(ctx.Request.Context(), &review, userID); err != nil {
		log.Printf("Error creating review: %v", err)
		// Handle specific errors from service, e.g., already reviewed, not registered
		if err.Error() == "you have already reviewed this event" || err.Error() == "user not registered for this event, cannot review" || err.Error() == "event not found" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Review created successfully", "review": review})
}

func (c *ReviewController) GetReviewsForEvent(ctx *gin.Context) {
	eventIDStr := ctx.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	reviews, err := c.reviewService.GetReviewsForEvent(ctx.Request.Context(), eventID)
	if err != nil {
		log.Printf("Error getting reviews for event %d: %v", eventID, err)
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "event not found" { // Check for specific "not found" cases
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found or no reviews available"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reviews"})
		}
		return
	}

	if len(reviews) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No reviews found for this event", "reviews": []model.Review{}})
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}
