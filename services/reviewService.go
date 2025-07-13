package services

import (
	"context"
	"errors"
	"fmt" // Added for fmt.Errorf
	"go-rest-api/model"
	"go-rest-api/repository"
	"log"

	"github.com/google/uuid"
)

type ReviewService interface {
	CreateReview(ctx context.Context, review *model.Review, userID uuid.UUID) error
	GetReviewsForEvent(ctx context.Context, eventID uuid.UUID) ([]model.Review, error)
	// CheckIfUserRegisteredForEvent(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (bool, error)
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	eventRepo  repository.EventRepository // To check if user is registered for the event
}

func NewReviewService(reviewRepo repository.ReviewRepository, eventRepo repository.EventRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		eventRepo:  eventRepo,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, review *model.Review, userID uuid.UUID) error {
	// Validate event exists
	_, err := s.eventRepo.GetEventById(ctx, review.EventID)
	if err != nil {
		log.Printf("Error finding event %d for review: %v", review.EventID, err)
		return errors.New("event not found")
	}

	// Check if the user has already reviewed this event
	existingReview, err := s.reviewRepo.GetReviewByEventAndUser(ctx, review.EventID, userID)
	if err != nil {
		// This is an actual error during DB query, not "not found"
		log.Printf("Error checking existing review for event %d by user %d: %v", review.EventID, userID, err)
		return errors.New("failed to check for existing reviews")
	}
	if existingReview != nil {
		return errors.New("you have already reviewed this event")
	}

	review.UserID = userID // Ensure the review is associated with the authenticated user
	err = s.reviewRepo.SaveReview(ctx, review)
	if err != nil {
		return err
	}

	// After saving a new review, recalculate and update the event's average rating
	go func() { // Run in a goroutine so it doesn't block the response
		err := s.recalculateAndUpdateAverageRating(context.Background(), review.EventID)
		if err != nil {
			log.Printf("Error updating average rating for event %d after new review: %v", review.EventID, err)
		}
	}()

	return nil
}

func (s *reviewService) GetReviewsForEvent(ctx context.Context, eventID uuid.UUID) ([]model.Review, error) {
	// Validate event exists (optional, as above)
	_, err := s.eventRepo.GetEventById(ctx, eventID)
	if err != nil {
		log.Printf("Error finding event %d when fetching reviews: %v", eventID, err)
		return nil, errors.New("event not found")
	}
	return s.reviewRepo.GetReviewsByEventID(ctx, eventID)
}

func (s *reviewService) recalculateAndUpdateAverageRating(ctx context.Context, eventID uuid.UUID) error {
	reviews, err := s.reviewRepo.GetReviewsByEventID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get reviews for recalculating average rating: %w", err)
	}

	// Use a new background context for this self-contained task.
	// This is crucial because the original context from the HTTP request
	// might be cancelled if the client disconnects or the request times out,
	// but still want this background update to complete if possible.
	updateCtx := context.Background()

	if len(reviews) == 0 {
		return s.eventRepo.UpdateAverageRating(updateCtx, eventID, 0) // Set to 0 if no reviews
	}

	var totalRating int
	for _, r := range reviews {
		totalRating += r.Rating
	}
	averageRating := float64(totalRating) / float64(len(reviews))

	return s.eventRepo.UpdateAverageRating(updateCtx, eventID, averageRating)
}

// TODO: Implement this method if you want to check user registration status before allowing review creation.
