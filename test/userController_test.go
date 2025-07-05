package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ctrl "go-rest-api/controllers"
	"go-rest-api/model"
	"go-rest-api/response"
)

// MockUserService is a mock type for the UserService type
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) ValidateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful registration", func(t *testing.T) {
		mockService := new(MockUserService)
		userController := ctrl.NewUserController(mockService)

		user := model.User{Email: "test@example.com", Password: "password123"}
		mockService.On("CreateUser", mock.AnythingOfType("*gin.Context"), &user).Return(nil)

		router := gin.New()
		router.POST("/register", userController.RegisterUser)

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var res gin.H
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.NotNil(t, res["user"])
		mockService.AssertExpectations(t)
	})

	t.Run("Email already registered", func(t *testing.T) {
		mockService := new(MockUserService)
		userController := ctrl.NewUserController(mockService)

		user := model.User{Email: "test@example.com", Password: "password123"}
		mockService.On("CreateUser", mock.AnythingOfType("*gin.Context"), &user).Return(errors.New("email already registered"))

		router := gin.New()
		router.POST("/register", userController.RegisterUser)

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var res gin.H
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, "email already registered", res["error"])
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid input", func(t *testing.T) {
		mockService := new(MockUserService)
		userController := ctrl.NewUserController(mockService)

		router := gin.New()
		router.POST("/register", userController.RegisterUser)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful login", func(t *testing.T) {
		mockService := new(MockUserService)
		userController := ctrl.NewUserController(mockService)

		user := model.User{Email: "test@example.com", Password: "password123"}
		mockService.On("ValidateUser", mock.AnythingOfType("*gin.Context"), &user).
			Run(func(args mock.Arguments) {
				arg := args.Get(1).(*model.User)
				arg.Id = 1
				arg.Role = "user"
			}).
			Return(nil)

		router := gin.New()
		router.POST("/login", userController.LoginUser)

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res gin.H
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.NotNil(t, res["token"])
		assert.NotNil(t, res["user"])

		userRes := res["user"].(map[string]interface{})
		assert.Equal(t, "test@example.com", userRes["email"])
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		mockService := new(MockUserService)
		userController := ctrl.NewUserController(mockService)

		user := model.User{Email: "wrong@example.com", Password: "wrongpassword"}
		mockService.On("ValidateUser", mock.AnythingOfType("*gin.Context"), &user).Return(errors.New("invalid credentials"))

		router := gin.New()
		router.POST("/login", userController.LoginUser)

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var res gin.H
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, "Invalid credentials", res["error"])
		mockService.AssertExpectations(t)
	})
}

func TestGetAllUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful retrieval", func(t *testing.T) {
		mockService := new(MockUserService)
		userController := ctrl.NewUserController(mockService)

		expectedUsers := []model.User{
			{Id: 1, Email: "user1@example.com", Role: "user"},
			{Id: 2, Email: "user2@example.com", Role: "admin"},
		}
		mockService.On("GetAllUsers", mock.AnythingOfType("*gin.Context")).Return(expectedUsers, nil)

		router := gin.New()
		router.GET("/users", userController.GetAllUser)

		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res gin.H
		json.Unmarshal(w.Body.Bytes(), &res)

		usersData := res["users"].([]interface{})
		assert.Len(t, usersData, 2)

		var userResponses []response.UserResponse
		json.Unmarshal(w.Body.Bytes(), &gin.H{"users": &userResponses})

		// This is a workaround to get the data into the struct slice
		data, _ := json.Marshal(res["users"])
		json.Unmarshal(data, &userResponses)

		assert.Equal(t, "user1@example.com", userResponses[0].Email)
		assert.Equal(t, "user2@example.com", userResponses[1].Email)
		mockService.AssertExpectations(t)
	})
}
