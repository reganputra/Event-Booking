package controllers

import (
	"database/sql"
	"errors"
	"go-rest-api/model"
	"go-rest-api/services"
	"go-rest-api/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	userID := userIDVal.(uuid.UUID)

	var event model.Event
	err := ctx.ShouldBindJSON(&event)
	if err != nil {
		validationErrors := utils.GetValidationErrors(err)
		if validationErrors != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
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

func (c *EventController) SearchEvents(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	startDate := ctx.Query("startDate") // Expected format: YYYY-MM-DD
	endDate := ctx.Query("endDate")     // Expected format: YYYY-MM-DD

	events, err := c.eventService.GetEventsByCriteria(ctx, keyword, startDate, endDate)
	if err != nil {
		log.Printf("Error searching events: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search events"})
		return
	}

	if len(events) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No events found matching your criteria", "events": []model.Event{}})
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
	eventID, err := uuid.Parse(id)
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
	userID := userIDVal.(uuid.UUID)
	userRoleVal, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context. Authentication issue."})
		return
	}
	userRole := userRoleVal.(string)

	id := ctx.Param("id")
	eventID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var event model.Event
	err = ctx.ShouldBindJSON(&event)
	if err != nil {
		validationErrors := utils.GetValidationErrors(err)
		if validationErrors != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
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
	userID := userIDVal.(uuid.UUID)
	userRoleVal, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context. Authentication issue."})
		return
	}
	userRole := userRoleVal.(string)

	id := ctx.Param("id")
	eventID, err := uuid.Parse(id)
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
	userID := userIDVal.(uuid.UUID)

	id := ctx.Param("id")
	eventID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = c.eventService.RegisterForEvent(ctx, eventID, userID)
	if err != nil {
		// Check for specific errors from the service, like "event is full, user added to waitlist"
		if err.Error() == "event is full, user added to waitlist" {
			ctx.JSON(http.StatusAccepted, gin.H{"message": err.Error()}) // 202 Accepted might be suitable
		} else if errors.Is(err, services.ErrAlreadyRegistered) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if errors.Is(err, services.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Error registering for event %d by user %d: %v", eventID, userID, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register for event"})
		}
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
	userID := userIDVal.(uuid.UUID)
	id := ctx.Param("id")
	eventID, err := uuid.Parse(id)
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
	userID := userIDVal.(uuid.UUID)

	events, err := c.eventService.GetRegisteredEvents(ctx, userID)
	if err != nil {
		log.Printf("Error getting registered events for user %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve registered events"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
