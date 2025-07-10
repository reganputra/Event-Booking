package controllers

import (
	"errors"
	"go-rest-api/model"
	"go-rest-api/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WaitlistController struct {
	waitlistService services.WaitlistService
	eventService    services.EventService
}

func NewWaitlistController(waitlistService services.WaitlistService, eventService services.EventService) *WaitlistController {
	return &WaitlistController{
		waitlistService: waitlistService,
		eventService:    eventService,
	}
}

// Join the waitlist for an event
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

// Leave the waitlist for an event
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

// Get the waitlist for an event (admin/owner only)
func (c *WaitlistController) GetWaitlistForEvent(ctx *gin.Context) {
	// Get current user ID and role from context
	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDVal.(int64)

	userRoleVal, exists := ctx.Get("role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
		return
	}
	userRole := userRoleVal.(string)

	eventIDStr := ctx.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	// Fetch the event to get the owner's ID
	event, err := c.eventService.GetEventByID(ctx.Request.Context(), eventID)
	if err != nil {
		log.Printf("Error fetching event %d for authorization check: %v", eventID, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Authorization check: user must be admin OR the event owner
	isAdmin := userRole == "admin"
	isOwner := event.UserIds == userID

	if !isAdmin && !isOwner {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied: you must be an admin or the event owner to view the waitlist"})
		return
	}

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
