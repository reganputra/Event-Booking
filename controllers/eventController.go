package controllers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-rest-api/model"
	"go-rest-api/utils"
	"net/http"
	"strconv"
)

func GetAllEvents(c *gin.Context) {
	events, err := model.GetAllEvents(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func GetEventsById(c *gin.Context) {

	id := c.Param("id")
	eventId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	event, err := model.GetEventById(c, eventId)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		}
		return
	}
	c.JSON(http.StatusOK, event)
}

func CreateEvent(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	userId, err := utils.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var createEvent model.Event
	err = c.ShouldBindJSON(&createEvent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createEvent.UserIds = userId
	err = createEvent.Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully!", "event": createEvent})
}

func UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	eventId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	_, err = model.GetEventById(c, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	var updateEvent model.Event
	err = c.ShouldBindJSON(&updateEvent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not update event"})
		return
	}
	updateEvent.Id = eventId
	err = updateEvent.UpdateEvent(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!", "event": updateEvent})

}

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	eventId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	event, err := model.GetEventById(c, eventId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		}
		return
	}
	err = event.DeleteEvent(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!"})
}
