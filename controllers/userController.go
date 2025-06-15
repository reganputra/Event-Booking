package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/model"
	"go-rest-api/response"
	"go-rest-api/services"
	"go-rest-api/utils"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) RegisterUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = u.userService.CreateUser(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	userResponse := response.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}

	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
}

func (u *UserController) LoginUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = u.userService.ValidateUser(c, &user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	userResponse := response.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}

	token, err := utils.GenerateToken(user.Email, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse, "token": token})
}
