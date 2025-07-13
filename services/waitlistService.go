package services

import (
	"context"
	"errors"
	"fmt"
	"go-rest-api/model"
	"go-rest-api/repository"
	"log"

	"github.com/google/uuid"
)

var ErrEventNotFull = errors.New("event is not full, cannot join waitlist")
var ErrAlreadyRegistered = errors.New("user is already registered for this event")
var ErrAlreadyOnWaitlist = errors.New("user is already on the waitlist for this event")
var ErrEventNotFound = errors.New("event not found")
var ErrUserNotOnWaitlist = errors.New("user is not on the waitlist for this event")
var ErrWaitlistNotEnabled = errors.New("waitlist not enabled for this event (capacity is 0 or not set)")

type WaitlistService interface {
	JoinWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (*model.WaitlistEntry, error)
	LeaveWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) error
	GetWaitlistForEvent(ctx context.Context, eventID uuid.UUID) ([]model.WaitlistEntry, error)
	ProcessNextOnWaitlist(ctx context.Context, eventID uuid.UUID) (*model.User, error)
}

type waitlistService struct {
	waitlistRepo repository.WaitlistRepository
	eventRepo    repository.EventRepository
	userRepo     repository.UserRepository // For fetching user details for notification (future)
	// notificationService NotificationService // For actual notifications (future)
}

func NewWaitlistService(
	waitlistRepo repository.WaitlistRepository,
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
	// notificationService NotificationService,
) WaitlistService {
	return &waitlistService{
		waitlistRepo: waitlistRepo,
		eventRepo:    eventRepo,
		userRepo:     userRepo,
		// notificationService: notificationService,
	}
}

func (s *waitlistService) JoinWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (*model.WaitlistEntry, error) {
	event, err := s.eventRepo.GetEventById(ctx, eventID)
	if err != nil {
		log.Printf("Error fetching event %d for waitlist join: %v", eventID, err)
		return nil, ErrEventNotFound
	}

	if event.Capacity <= 0 {
		return nil, ErrWaitlistNotEnabled
	}

	registeredCount, err := s.eventRepo.GetRegistrationCount(ctx, eventID)
	if err != nil {
		log.Printf("Error fetching registration count for event %d: %v", eventID, err)
		return nil, fmt.Errorf("could not verify event registration count: %w", err)
	}

	if registeredCount < event.Capacity {
		return nil, ErrEventNotFull
	}

	isRegistered, err := s.eventRepo.IsUserRegistered(ctx, eventID, userID)
	if err != nil {
		log.Printf("Error checking if user %d is registered for event %d: %v", userID, eventID, err)
		return nil, fmt.Errorf("could not verify event registration status: %w", err)
	}
	if isRegistered {
		return nil, ErrAlreadyRegistered
	}

	isOnWaitlist, err := s.waitlistRepo.IsUserOnWaitlist(ctx, eventID, userID)
	if err != nil {
		log.Printf("Error checking if user %d is on waitlist for event %d: %v", userID, eventID, err)
		return nil, fmt.Errorf("could not verify waitlist status: %w", err)
	}
	if isOnWaitlist {
		return nil, ErrAlreadyOnWaitlist
	}

	return s.waitlistRepo.AddUserToWaitlist(ctx, eventID, userID)
}

func (s *waitlistService) LeaveWaitlist(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) error {
	// Check if event exists (optional, FK constraint might cover it)
	_, err := s.eventRepo.GetEventById(ctx, eventID)
	if err != nil {
		return ErrEventNotFound
	}

	isOnWaitlist, err := s.waitlistRepo.IsUserOnWaitlist(ctx, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to check waitlist status: %w", err)
	}
	if !isOnWaitlist {
		return ErrUserNotOnWaitlist
	}

	return s.waitlistRepo.RemoveUserFromWaitlist(ctx, eventID, userID)
}

func (s *waitlistService) GetWaitlistForEvent(ctx context.Context, eventID uuid.UUID) ([]model.WaitlistEntry, error) {
	_, err := s.eventRepo.GetEventById(ctx, eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}
	return s.waitlistRepo.GetWaitlistForEvent(ctx, eventID)
}

// ProcessNextOnWaitlist is called when a spot opens up (e.g., someone cancels registration).
// This is a simplified version. A real system might:
// - Actually register the user.
// - Send a notification with a time limit to register.
// - Handle cases where the next user is no longer interested.
func (s *waitlistService) ProcessNextOnWaitlist(ctx context.Context, eventID uuid.UUID) (*model.User, error) {
	log.Printf("Processing next on waitlist for event ID %s", eventID)
	nextEntry, err := s.waitlistRepo.GetNextUserFromWaitlist(ctx, eventID)
	if err != nil {
		log.Printf("Error getting next user from waitlist for event %d: %v", eventID, err)
		return nil, fmt.Errorf("failed to get next user from waitlist: %w", err)
	}

	if nextEntry == nil {
		log.Printf("No users on waitlist for event %d", eventID)
		return nil, nil // No one to process
	}

	// "Promote" the user: For now, this means removing them from waitlist and attempting to register them.
	// In a more complex system, you'd notify them.

	// Attempt to register the user directly.
	// This assumes the spot is indeed free. A lock might be needed in high concurrency.
	err = s.eventRepo.RegisterEvent(ctx, eventID, nextEntry.UserID)
	if err != nil {
		// This could happen if, by some race condition, the event filled up again,
		// or the user got registered through another means.
		log.Printf("Failed to auto-register user %d from waitlist for event %d: %v", nextEntry.UserID, eventID, err)
		// Depending on policy, you might leave them on waitlist or try the next person.
		// For now, we'll just log and indicate no user was successfully promoted.
		return nil, fmt.Errorf("failed to register user from waitlist: %w", err)
	}

	// If registration was successful, remove them from the waitlist.
	err = s.waitlistRepo.RemoveUserFromWaitlist(ctx, eventID, nextEntry.UserID)
	if err != nil {
		// This is problematic: they are registered but still on waitlist. Needs careful handling.
		log.Printf("CRITICAL: User %d registered from waitlist for event %d but failed to remove from waitlist: %v", nextEntry.UserID, eventID, err)
		// Return the user since they were registered, but this state is inconsistent.
	}

	log.Printf("User %d successfully registered from waitlist for event %d.", nextEntry.UserID, eventID)

	promotedUser, userErr := s.userRepo.GetById(ctx, nextEntry.UserID) // Corrected method call
	if userErr != nil {
		log.Printf("Failed to fetch details for promoted user %d: %v", nextEntry.UserID, userErr)
		// User was promoted, but we couldn't fetch their details.
		// Return a placeholder or handle error as per requirements.
		// For now, we'll return nil for the user object but no error for the promotion itself.
		return nil, nil
	}

	// Here you would typically send a notification.
	// s.notificationService.SendWaitlistPromotionNotification(promotedUser, event)

	return promotedUser, nil
}
