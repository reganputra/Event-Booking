package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/model"
	"net/http"
	"strconv"
)

func RegisterForEvent(c *gin.Context) {
	userId := c.GetInt64("userId")
	id := c.Param("id")
	eventId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	event, err := model.GetEventById(c, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	err = event.RegisterEvent(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register for event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered for the event", "event": event})
}

func CancelForEvent(c *gin.Context) {
	userId := c.GetInt64("userId")
	id := c.Param("id")
	eventId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	event, err := model.GetEventById(c, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	err = event.CancelRegistration(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel event registration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully cancelled event registration"})
}
