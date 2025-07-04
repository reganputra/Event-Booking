package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-rest-api/controllers"
	"go-rest-api/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReviewService is a mock type for the ReviewService type
type MockReviewService struct {
	mock.Mock
}

func (m *MockReviewService) CreateReview(ctx context.Context, review *model.Review, userID int64) error {
	args := m.Called(ctx, review, userID)
	// Modify the review object passed in, similar to how the actual service would (e.g., setting ID, CreatedAt)
	if args.Error(0) == nil {
		review.Id = 1 // Simulate ID generation
		review.CreatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockReviewService) GetReviewsForEvent(ctx context.Context, eventID int64) ([]model.Review, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Review), args.Error(1)
}

func setupReviewRouter(reviewService *MockReviewService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	reviewController := controllers.NewReviewController(reviewService)

	// Middleware to mock authentication (set user ID)
	authMiddleware := func(c *gin.Context) {
		c.Set("userId", int64(1)) // Mock user ID
		c.Next()
	}

	router.POST("/events/:id/reviews", authMiddleware, reviewController.CreateReview)
	router.GET("/events/:id/reviews", reviewController.GetReviewsForEvent)
	return router
}

func TestCreateReview(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupReviewRouter(mockService)

	t.Run("Successful review creation", func(t *testing.T) {
		reviewInput := model.Review{EventID: 1, Rating: 5, Comment: "Great event!"}
		mockService.On("CreateReview", mock.Anything, mock.AnythingOfType("*model.Review"), int64(1)).Return(nil).Once()

		body, _ := json.Marshal(reviewInput)
		req, _ := http.NewRequest(http.MethodPost, "/events/1/reviews", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Review created successfully", response["message"])
		reviewData, ok := response["review"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(reviewInput.Rating), reviewData["rating"]) // JSON numbers are float64
		assert.Equal(t, reviewInput.Comment, reviewData["comment"])
		mockService.AssertExpectations(t)
	})

	t.Run("Review creation - already reviewed", func(t *testing.T) {
		reviewInput := model.Review{EventID: 1, Rating: 3, Comment: "Trying again"}
		mockService.On("CreateReview", mock.Anything, mock.AnythingOfType("*model.Review"), int64(1)).Return(errors.New("you have already reviewed this event")).Once()

		body, _ := json.Marshal(reviewInput)
		req, _ := http.NewRequest(http.MethodPost, "/events/1/reviews", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "you have already reviewed this event", response["error"])
		mockService.AssertExpectations(t)
	})

	t.Run("Review creation - event not found", func(t *testing.T) {
		reviewInput := model.Review{EventID: 999, Rating: 5, Comment: "For non-existent event"}
		mockService.On("CreateReview", mock.Anything, mock.AnythingOfType("*model.Review"), int64(1)).Return(errors.New("event not found")).Once()

		body, _ := json.Marshal(reviewInput)
		req, _ := http.NewRequest(http.MethodPost, "/events/999/reviews", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code) // Based on controller logic for this error
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "event not found", response["error"])
		mockService.AssertExpectations(t)
	})

	t.Run("Review creation - invalid input (rating out of bounds)", func(t *testing.T) {
		// This test doesn't need mockService as validation is by Gin binding
		reviewInput := gin.H{"event_id": 1, "rating": 0, "comment": "Bad rating"} // Rating 0

		body, _ := json.Marshal(reviewInput)
		req, _ := http.NewRequest(http.MethodPost, "/events/1/reviews", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// No need to assert mockService.AssertExpectations(t) as the service shouldn't be called
	})

}

func TestGetReviewsForEvent(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupReviewRouter(mockService)

	t.Run("Successful retrieval of reviews", func(t *testing.T) {
		expectedReviews := []model.Review{
			{Id: 1, EventID: 1, UserID: 1, Rating: 5, Comment: "Excellent!", CreatedAt: time.Now()},
			{Id: 2, EventID: 1, UserID: 2, Rating: 4, Comment: "Very good.", CreatedAt: time.Now().Add(-time.Hour)},
		}
		mockService.On("GetReviewsForEvent", mock.Anything, int64(1)).Return(expectedReviews, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/1/reviews", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseReviews []model.Review
		err := json.Unmarshal(w.Body.Bytes(), &responseReviews)
		assert.NoError(t, err)
		assert.Len(t, responseReviews, 2)
		// For more detailed comparison, especially with time.Time, you might need custom comparators
		// or iterate and assert fields. For simplicity, we check length and IDs if available.
		assert.Equal(t, expectedReviews[0].Id, responseReviews[0].Id)
		assert.Equal(t, expectedReviews[1].Id, responseReviews[1].Id)
		mockService.AssertExpectations(t)
	})

	t.Run("No reviews found for event", func(t *testing.T) {
		mockService.On("GetReviewsForEvent", mock.Anything, int64(2)).Return([]model.Review{}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/2/reviews", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "No reviews found for this event", response["message"])
		assert.Len(t, response["reviews"], 0)
		mockService.AssertExpectations(t)
	})

	t.Run("Event not found when getting reviews", func(t *testing.T) {
		mockService.On("GetReviewsForEvent", mock.Anything, int64(999)).Return(nil, errors.New("event not found")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/999/reviews", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code) // Controller maps "event not found" to 404
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Event not found or no reviews available", response["error"])
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error when getting reviews", func(t *testing.T) {
		mockService.On("GetReviewsForEvent", mock.Anything, int64(3)).Return(nil, fmt.Errorf("some database error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/3/reviews", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve reviews", response["error"])
		mockService.AssertExpectations(t)
	})
}
