package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/model"
	"go-rest-api/response"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	err = user.CreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userResponse := response.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!", "user": userResponse})
}

func LoginUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	err = user.ValidateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userResponse := response.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully!", "user": userResponse})
}
