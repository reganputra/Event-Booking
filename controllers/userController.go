package controllers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-rest-api/model"
	"go-rest-api/response"
	"go-rest-api/services"
	"go-rest-api/utils"
	"net/http"
	"strconv"
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
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	userResponse := response.UserResponse{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
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
		Role:  user.Role,
	}

	token, err := utils.GenerateToken(user.Email, user.Id, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse, "token": token})
}

func (u *UserController) GetAllUser(c *gin.Context) {
	users, err := u.userService.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	var userResponses []response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, response.UserResponse{
			Id:    user.Id,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": userResponses})
}

func (u *UserController) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	user, err := u.userService.GetUserByID(c, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}
	userResponse := response.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, gin.H{"user": userResponse})
}

func (u *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var user model.User
	err = c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user.Id = id

	err = u.userService.UpdateUser(c, &user)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else if err.Error() == "invalid role" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		}
		return
	}

	userResponse := response.UserResponse{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}
	c.JSON(http.StatusOK, gin.H{"user": userResponse})
}

func (u *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = u.userService.DeleteUser(c, id)
	if err != nil {
		if err.Error() == "user not found" || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted successfully"})
}
