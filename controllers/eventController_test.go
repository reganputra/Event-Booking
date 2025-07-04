package controllers_test

import (
	"bytes"
	"context" // Import the standard context package
	"encoding/json"
	"go-rest-api/model"
	// "go-rest-api/services" // No longer directly used due to mock interface
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

import ctrl "go-rest-api/controllers" // Alias import for controllers

// MockEventService is a mock type for the EventService type
type MockEventService struct {
	mock.Mock
}

// CreateEvent is a mock method
func (m *MockEventService) CreateEvent(ctx context.Context, event *model.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// GetAllEvents is a mock method
func (m *MockEventService) GetAllEvents(ctx context.Context) ([]model.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Event), args.Error(1)
}

// GetEventByID is a mock method
func (m *MockEventService) GetEventByID(ctx context.Context, id int64) (*model.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Event), args.Error(1)
}

// GetEventsByCategory is a mock method
func (m *MockEventService) GetEventsByCategory(ctx context.Context, category string) ([]model.Event, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]model.Event), args.Error(1)
}

// GetEventsByCriteria is a mock method
func (m *MockEventService) GetEventsByCriteria(ctx context.Context, keyword string, startDate string, endDate string) ([]model.Event, error) {
	args := m.Called(ctx, keyword, startDate, endDate)
	return args.Get(0).([]model.Event), args.Error(1)
}

// UpdateEvent is a mock method
func (m *MockEventService) UpdateEvent(ctx context.Context, event *model.Event, userID int64, userRole string) error {
	args := m.Called(ctx, event, userID, userRole)
	return args.Error(0)
}

// DeleteEvent is a mock method
func (m *MockEventService) DeleteEvent(ctx context.Context, id int64, userID int64, userRole string) error {
	args := m.Called(ctx, id, userID, userRole)
	return args.Error(0)
}

// RegisterForEvent is a mock method
func (m *MockEventService) RegisterForEvent(ctx context.Context, eventID, userID int64) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

// CancelEventRegistration is a mock method
func (m *MockEventService) CancelEventRegistration(ctx context.Context, eventID, userID int64) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

// GetRegisteredEvents is a mock method
func (m *MockEventService) GetRegisteredEvents(ctx context.Context, userID int64) ([]model.Event, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Event), args.Error(1)
}

// Context is an alias for gin.Context for brevity in mock calls
// type Context = *gin.Context // This was causing an import cycle, using services.Context directly

func TestSearchEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful search with keyword", func(t *testing.T) {
		mockService := new(MockEventService)
		// Use the aliased import for NewEventController
		eventController := ctrl.NewEventController(mockService)

		expectedEvents := []model.Event{
			{Id: 1, Name: "Tech Conference", Description: "A conference about technology", Location: "Online", Date: time.Now(), Category: "Tech", UserIds: 1},
		}
		// Use mock.Anything for the context argument
		mockService.On("GetEventsByCriteria", mock.AnythingOfType("*gin.Context"), "Tech", "", "").Return(expectedEvents, nil)


		router := gin.New()
		router.GET("/events/search", eventController.SearchEvents)

		req, _ := http.NewRequest(http.MethodGet, "/events/search?keyword=Tech", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseEvents []model.Event
		err := json.Unmarshal(w.Body.Bytes(), &responseEvents)
		assert.NoError(t, err)
		// Comparing time.Time objects can be tricky due to monotonic clock.
		// For simplicity here, if IDs match, we assume it's good enough for this test structure.
		// A more robust check would compare individual fields or use assert.WithinDuration for time.
		assert.Equal(t, len(expectedEvents), len(responseEvents))
		if len(expectedEvents) > 0 && len(responseEvents) > 0 {
			assert.Equal(t, expectedEvents[0].Id, responseEvents[0].Id)
		}
		mockService.AssertExpectations(t)
	})

	t.Run("Successful search with date range", func(t *testing.T) {
		mockService := new(MockEventService)
		eventController := ctrl.NewEventController(mockService)

		startDate := "2024-01-01"
		endDate := "2024-01-31"
		// Correctly initialize time.Time for expectedEvents to avoid issues with time zone or monotonic clock parts
		loc, _ := time.LoadLocation("UTC") // Or use time.Local if your app logic implies local time
		expectedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, loc)
		expectedEvents := []model.Event{
			{Id: 2, Name: "Workshop", Description: "A workshop event", Location: "Local Hub", Date: expectedDate, Category: "Education", UserIds: 2},
		}
		mockService.On("GetEventsByCriteria", mock.AnythingOfType("*gin.Context"), "", startDate, endDate).Return(expectedEvents, nil)


		router := gin.New()
		router.GET("/events/search", eventController.SearchEvents)

		req, _ := http.NewRequest(http.MethodGet, "/events/search?startDate=2024-01-01&endDate=2024-01-31", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseEvents []model.Event
		err := json.Unmarshal(w.Body.Bytes(), &responseEvents)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedEvents), len(responseEvents))
		if len(expectedEvents) > 0 && len(responseEvents) > 0 {
			assert.Equal(t, expectedEvents[0].Id, responseEvents[0].Id)
			// For time comparison, it's often better to compare formatted strings or Unix timestamps
			// if direct equality fails due to subtle differences (like monotonic clock).
			assert.True(t, expectedEvents[0].Date.Equal(responseEvents[0].Date), "Dates should be equal")
		}
		mockService.AssertExpectations(t)
	})

	t.Run("No events found", func(t *testing.T) {
		mockService := new(MockEventService)
		eventController := ctrl.NewEventController(mockService)

		mockService.On("GetEventsByCriteria", mock.AnythingOfType("*gin.Context"), "NonExistent", "", "").Return([]model.Event{}, nil)


		router := gin.New()
		router.GET("/events/search", eventController.SearchEvents)

		req, _ := http.NewRequest(http.MethodGet, "/events/search?keyword=NonExistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "No events found matching your criteria", response["message"])
		assert.NotNil(t, response["events"]) // Ensure events field is present, even if empty
		actualEvents, ok := response["events"].([]interface{})
		assert.True(t, ok, "events field should be a slice")
		assert.Len(t, actualEvents, 0, "events slice should be empty")
		mockService.AssertExpectations(t)
	})

}

// --- Helper function to create a request with a JSON body ---
func newJsonRequest(method, url string, body interface{}) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
}

// --- Mocking context for user ID and role ---
func getMockedContext(req *http.Request, w *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	c, r := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userId", int64(1))    // Mock user ID
	c.Set("userRole", "user") // Mock user role
	return c, r
}
