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

type EventController struct {
	eventService services.EventService
}

func NewEventController(eventService services.EventService) *EventController {
	return &EventController{eventService: eventService}
}

func (c *EventController) CreateEvent(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication issue."})
		return
	}
	userID := userIDVal.(int64)

	var event model.Event
	err := ctx.ShouldBindJSON(&event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.UserIds = userID
	err = c.eventService.CreateEvent(ctx, &event)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Event created successfully!", "event": event})
}

func (c *EventController) GetAllEvents(ctx *gin.Context) {
	events, err := c.eventService.GetAllEvents(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}
	ctx.JSON(http.StatusOK, events)
}

func (c *EventController) GetEventsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	if category == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category parameter is required"})
		return
	}

	events, err := c.eventService.GetEventsByCategory(ctx, category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events by category"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func (c *EventController) GetEventByID(ctx *gin.Context) {
	id := ctx.Param("id")
	eventID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	event, err := c.eventService.GetEventByID(ctx, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		}
		return
	}
	ctx.JSON(http.StatusOK, event)
}

func (c *EventController) UpdateEvent(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication issue."})
		return
	}
	userID := userIDVal.(int64)
	userRoleVal, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context. Authentication issue."})
		return
	}
	userRole := userRoleVal.(string)

	id := ctx.Param("id")
	eventID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var event model.Event
	err = ctx.ShouldBindJSON(&event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not update event"})
		return
	}

	event.Id = eventID
	err = c.eventService.UpdateEvent(ctx, &event, userID, userRole)
	if err != nil {
		if err.Error() == "unauthorized: you don't have permission to update this event" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!", "event": event})
}

func (c *EventController) DeleteEvent(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication issue."})
		return
	}
	userID := userIDVal.(int64)
	userRoleVal, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context. Authentication issue."})
		return
	}
	userRole := userRoleVal.(string)

	id := ctx.Param("id")
	eventID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = c.eventService.DeleteEvent(ctx, eventID, userID, userRole)
	if err != nil {
		if err.Error() == "unauthorized: you don't have permission to delete this event" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!"})
}

func (c *EventController) RegisterForEvent(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication issue."})
		return
	}
	userID := userIDVal.(int64)

	id := ctx.Param("id")
	eventID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = c.eventService.RegisterForEvent(ctx, eventID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register for event"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully registered for the event"})
}

func (c *EventController) CancelEventRegistration(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication issue."})
		return
	}
	userID := userIDVal.(int64)
	id := ctx.Param("id")
	eventID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = c.eventService.CancelEventRegistration(ctx, eventID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel event registration"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully cancelled event registration"})
}

func (c *EventController) GetRegisteredEvents(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context. Authentication issue."})
		return
	}
	userID := userIDVal.(int64)

	events, err := c.eventService.GetRegisteredEvents(ctx, userID)
	if err != nil {
		log.Printf("Error getting registered events for user %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve registered events"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
