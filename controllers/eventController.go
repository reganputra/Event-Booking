package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/model"
	"net/http"
)

func GetAllEvents(c *gin.Context) {
	events := model.GetAllEvents(c)
	c.JSON(http.StatusOK, events)
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
