package controllers

import (
	"errors"
	"go-rest-api/model" // Added import for model
	"go-rest-api/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WaitlistController struct {
	waitlistService services.WaitlistService
}

func NewWaitlistController(waitlistService services.WaitlistService) *WaitlistController {
	return &WaitlistController{waitlistService: waitlistService}
}

// POST /events/:id/waitlist - Join the waitlist for an event
func (c *WaitlistController) JoinWaitlist(ctx *gin.Context) {
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

	entry, err := c.waitlistService.JoinWaitlist(ctx.Request.Context(), eventID, userID)
	if err != nil {
		log.Printf("Error joining waitlist for event %d by user %d: %v", eventID, userID, err)
		if errors.Is(err, services.ErrEventNotFull) ||
			errors.Is(err, services.ErrAlreadyRegistered) ||
			errors.Is(err, services.ErrAlreadyOnWaitlist) ||
			errors.Is(err, services.ErrEventNotFound) ||
		    errors.Is(err, services.ErrWaitlistNotEnabled) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join waitlist"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully joined the waitlist", "waitlist_entry": entry})
}

// DELETE /events/:id/waitlist - Leave the waitlist for an event
func (c *WaitlistController) LeaveWaitlist(ctx *gin.Context) {
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

	err = c.waitlistService.LeaveWaitlist(ctx.Request.Context(), eventID, userID)
	if err != nil {
		log.Printf("Error leaving waitlist for event %d by user %d: %v", eventID, userID, err)
		if errors.Is(err, services.ErrUserNotOnWaitlist) || errors.Is(err, services.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave waitlist"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully left the waitlist"})
}

// GET /events/:id/waitlist - Get the waitlist for an event (admin/owner only)
func (c *WaitlistController) GetWaitlistForEvent(ctx *gin.Context) {
	// Authentication and Authorization (e.g. checking if user is admin or event owner)
	// should be handled by middleware. For this example, we assume it's done.
	// We'll need to get the event owner's ID from the event itself if we want to allow owners.
	// For now, let's assume this is an admin-only endpoint or public for simplicity in this snippet.
	// The plan says "for event owner/admin", so this needs proper auth.

	eventIDStr := ctx.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	// TODO: Add authorization logic here:
	// 1. Get current user ID and role from ctx.
	// 2. Fetch event details to get event.UserID (owner).
	// 3. If current user is not admin AND current user ID is not event.UserID, return 403 Forbidden.

	entries, err := c.waitlistService.GetWaitlistForEvent(ctx.Request.Context(), eventID)
	if err != nil {
		log.Printf("Error getting waitlist for event %d: %v", eventID, err)
		if errors.Is(err, services.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve waitlist"})
		}
		return
	}

	if len(entries) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Waitlist is empty for this event", "waitlist": []model.WaitlistEntry{}})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"waitlist": entries})
}
