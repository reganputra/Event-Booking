package controllers_test

import (
	"context"
	"encoding/json"
	// "errors" // Unused import
	"fmt"
	"go-rest-api/controllers"
	"go-rest-api/model"
	"go-rest-api/services"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWaitlistService is a mock type for the WaitlistService type
type MockWaitlistService struct {
	mock.Mock
}

func (m *MockWaitlistService) JoinWaitlist(ctx context.Context, eventID int64, userID int64) (*model.WaitlistEntry, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WaitlistEntry), args.Error(1)
}

func (m *MockWaitlistService) LeaveWaitlist(ctx context.Context, eventID int64, userID int64) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

func (m *MockWaitlistService) GetWaitlistForEvent(ctx context.Context, eventID int64) ([]model.WaitlistEntry, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.WaitlistEntry), args.Error(1)
}

func (m *MockWaitlistService) ProcessNextOnWaitlist(ctx context.Context, eventID int64) (*model.User, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}


func setupWaitlistRouter(waitlistService services.WaitlistService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	waitlistController := controllers.NewWaitlistController(waitlistService)

	// Middleware to mock authentication (set user ID and role)
	authMiddleware := func(c *gin.Context) {
		c.Set("userId", int64(1)) // Mock user ID for non-admin
		// For admin-specific tests, this might need to be dynamic or set differently
		// For GetWaitlistForEvent, if testing admin access, role should be 'admin'
		// For Join/Leave, 'user' role is fine.
		// Let's assume role is 'user' for Join/Leave and 'admin' for Get for now.
		if c.Request.Method == "GET" && c.FullPath() == "/admin/events/:id/waitlist" {
			c.Set("userRole", "admin")
		} else {
			c.Set("userRole", "user")
		}
		c.Next()
	}

	// This matches how routes are defined in main.go
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(authMiddleware) // Apply to all user-level waitlist actions
	{
		protectedRoutes.POST("/events/:id/waitlist", waitlistController.JoinWaitlist)
		protectedRoutes.DELETE("/events/:id/waitlist", waitlistController.LeaveWaitlist)
	}

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(authMiddleware) // Ensures userId is set, role check is effectively done by path + mock
	{
		adminRoutes.GET("/events/:id/waitlist", waitlistController.GetWaitlistForEvent)
	}

	return router
}


func TestJoinWaitlist(t *testing.T) {
	mockService := new(MockWaitlistService)
	router := setupWaitlistRouter(mockService) // Use the correct type here

	t.Run("Successful join waitlist", func(t *testing.T) {
		eventID := int64(1)
		userID := int64(1) // From authMiddleware
		expectedEntry := &model.WaitlistEntry{Id: 100, EventID: eventID, UserID: userID, CreatedAt: time.Now()}

		mockService.On("JoinWaitlist", mock.Anything, eventID, userID).Return(expectedEntry, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/events/%d/waitlist", eventID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Successfully joined the waitlist", response["message"])
		entryData, _ := response["waitlist_entry"].(map[string]interface{})
		assert.Equal(t, float64(expectedEntry.Id), entryData["id"])
		mockService.AssertExpectations(t)
	})

	t.Run("Join waitlist - event not full", func(t *testing.T) {
		eventID := int64(2)
		userID := int64(1)
		mockService.On("JoinWaitlist", mock.Anything, eventID, userID).Return(nil, services.ErrEventNotFull).Once()

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/events/%d/waitlist", eventID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, services.ErrEventNotFull.Error(), response["error"])
		mockService.AssertExpectations(t)
	})

    t.Run("Join waitlist - already registered", func(t *testing.T) {
		eventID := int64(3)
		userID := int64(1)
		mockService.On("JoinWaitlist", mock.Anything, eventID, userID).Return(nil, services.ErrAlreadyRegistered).Once()

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/events/%d/waitlist", eventID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, services.ErrAlreadyRegistered.Error(), response["error"])
		mockService.AssertExpectations(t)
	})
}

func TestLeaveWaitlist(t *testing.T) {
	mockService := new(MockWaitlistService)
	router := setupWaitlistRouter(mockService)

	t.Run("Successful leave waitlist", func(t *testing.T) {
		eventID := int64(1)
		userID := int64(1)
		mockService.On("LeaveWaitlist", mock.Anything, eventID, userID).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/events/%d/waitlist", eventID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Successfully left the waitlist", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Leave waitlist - user not on waitlist", func(t *testing.T) {
		eventID := int64(2)
		userID := int64(1)
		mockService.On("LeaveWaitlist", mock.Anything, eventID, userID).Return(services.ErrUserNotOnWaitlist).Once()

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/events/%d/waitlist", eventID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, services.ErrUserNotOnWaitlist.Error(), response["error"])
		mockService.AssertExpectations(t)
	})
}

func TestGetWaitlistForEvent(t *testing.T) {
	mockService := new(MockWaitlistService)
	router := setupWaitlistRouter(mockService)

	t.Run("Successful get waitlist (admin)", func(t *testing.T) {
		eventID := int64(1)
		expectedEntries := []model.WaitlistEntry{
			{Id: 1, EventID: eventID, UserID: 10, CreatedAt: time.Now().Add(-2 * time.Hour)},
			{Id: 2, EventID: eventID, UserID: 11, CreatedAt: time.Now().Add(-1 * time.Hour)},
		}
		mockService.On("GetWaitlistForEvent", mock.Anything, eventID).Return(expectedEntries, nil).Once()

		// Note: The route for admin is /admin/events/:id/waitlist
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/admin/events/%d/waitlist", eventID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		waitlistData, ok := response["waitlist"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, waitlistData, 2)
		// Further checks on content if necessary
		mockService.AssertExpectations(t)
	})

	t.Run("Get waitlist - event not found (admin)", func(t *testing.T) {
		eventID := int64(999)
		mockService.On("GetWaitlistForEvent", mock.Anything, eventID).Return(nil, services.ErrEventNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/admin/events/%d/waitlist", eventID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, services.ErrEventNotFound.Error(), response["error"])
		mockService.AssertExpectations(t)
	})

    t.Run("Get waitlist - empty (admin)", func(t *testing.T) {
        eventID := int64(2)
        mockService.On("GetWaitlistForEvent", mock.Anything, eventID).Return([]model.WaitlistEntry{}, nil).Once()

        req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/admin/events/%d/waitlist", eventID), nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        var response gin.H
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Equal(t, "Waitlist is empty for this event", response["message"])
        waitlistData, _ := response["waitlist"].([]interface{})
        assert.Len(t, waitlistData, 0)
        mockService.AssertExpectations(t)
    })
}
