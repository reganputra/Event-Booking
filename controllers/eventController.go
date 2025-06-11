package controllers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-rest-api/model"
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
	var createEvent model.Event
	err := c.ShouldBindJSON(&createEvent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createEvent.Id = 1
	createEvent.UserIds = 1

	err = createEvent.Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully!", "event": createEvent})
}
