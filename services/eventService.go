package services

import (
	"context"
	"errors"
	"fmt" // Added import for fmt
	"go-rest-api/model"
	"go-rest-api/repository"
	"log" // Added import for log

	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	GetAllEvents(ctx context.Context) ([]model.Event, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
	GetEventsByCategory(ctx context.Context, category string) ([]model.Event, error)
	GetEventsByCriteria(ctx context.Context, keyword string, startDate string, endDate string) ([]model.Event, error)
	UpdateEvent(ctx context.Context, event *model.Event, userID uuid.UUID, userRole string) error
	DeleteEvent(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole string) error
	RegisterForEvent(ctx context.Context, eventID, userID uuid.UUID) error
	CancelEventRegistration(ctx context.Context, eventID, userID uuid.UUID) error
	GetRegisteredEvents(ctx context.Context, userID uuid.UUID) ([]model.Event, error)
}

type eventService struct {
	eventRepository repository.EventRepository
	waitlistService WaitlistService // Added to call ProcessNextOnWaitlist
}

func NewEventService(eventRepository repository.EventRepository, waitlistService WaitlistService) EventService {
	return &eventService{
		eventRepository: eventRepository,
		waitlistService: waitlistService,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, event *model.Event) error {
	// Default capacity to 0 if not provided or negative, unless binding already handles gte=0
	if event.Capacity < 0 {
		event.Capacity = 0
	}
	return s.eventRepository.Save(ctx, event)
}

func (s *eventService) GetAllEvents(ctx context.Context) ([]model.Event, error) {
	return s.eventRepository.GetAllEvents(ctx)
}

func (s *eventService) GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.eventRepository.GetEventById(ctx, id)
}

func (s *eventService) UpdateEvent(ctx context.Context, event *model.Event, userID uuid.UUID, userRole string) error {
	existingEvent, err := s.eventRepository.GetEventById(ctx, event.Id)
	if err != nil {
		return err
	}

	if existingEvent.UserIds != userID && userRole != "admin" {
		return errors.New("unauthorized: you don't have permission to update this event")
	}
	// Preserve existing capacity if not provided in update payload
	if event.Capacity == 0 && existingEvent.Capacity > 0 { // Check if capacity is being explicitly set to 0 or just omitted

	}

	return s.eventRepository.Update(ctx, event)
}

func (s *eventService) DeleteEvent(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole string) error {
	existingEvent, err := s.eventRepository.GetEventById(ctx, id)
	if err != nil {
		return err
	}

	if existingEvent.UserIds != userID && userRole != "admin" {
		return errors.New("unauthorized: you don't have permission to delete this event")
	}

	return s.eventRepository.DeleteEvent(ctx, id)
}

func (s *eventService) RegisterForEvent(ctx context.Context, eventID, userID uuid.UUID) error {
	event, err := s.eventRepository.GetEventById(ctx, eventID)
	if err != nil {
		return ErrEventNotFound // Use defined error
	}

	// Check if user is already registered
	isRegistered, err := s.eventRepository.IsUserRegistered(ctx, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to check registration status: %w", err)
	}
	if isRegistered {
		return ErrAlreadyRegistered // Use defined error
	}

	// Check capacity if it's set (event.Capacity > 0)
	if event.Capacity > 0 {
		registrationCount, err := s.eventRepository.GetRegistrationCount(ctx, eventID)
		if err != nil {
			return fmt.Errorf("failed to get registration count: %w", err)
		}
		if registrationCount >= event.Capacity {
			// Event is full, try adding to waitlist via WaitlistService
			log.Printf("Event %d is full. Attempting to add user %d to waitlist.", eventID, userID)
			_, wlErr := s.waitlistService.JoinWaitlist(ctx, eventID, userID)
			if wlErr != nil {
				log.Printf("Failed to add user %d to waitlist for event %d: %v", userID, eventID, wlErr)
				return fmt.Errorf("event is full and failed to join waitlist: %w", wlErr)
			}
			return errors.New("event is full, user added to waitlist") // Specific error/message
		}
	}

	return s.eventRepository.RegisterEvent(ctx, eventID, userID)
}

func (s *eventService) CancelEventRegistration(ctx context.Context, eventID, userID uuid.UUID) error {
	// Get event details before cancellation
	event, err := s.eventRepository.GetEventById(ctx, eventID)
	if err != nil {
		return err // Event not found
	}

	// Check if the event is currently at full capacity
	registrationCount, err := s.eventRepository.GetRegistrationCount(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get registration count: %w", err)
	}

	// Determine if the event is currently at full capacity
	isFull := event.Capacity > 0 && registrationCount >= event.Capacity

	// Cancel the registration
	err = s.eventRepository.CancelRegistration(ctx, eventID, userID)
	if err != nil {
		return err // Failed to cancel or user wasn't registered
	}

	// Only process the waitlist if the event was at full capacity before cancellation
	if isFull {
		// Run in a goroutine to avoid blocking the cancellation response
		go func() {
			log.Printf("Event %d was at full capacity. Processing waitlist after cancellation.", eventID)
			promotedUser, err := s.waitlistService.ProcessNextOnWaitlist(context.Background(), eventID)
			if err != nil {
				log.Printf("Error processing waitlist for event %d after cancellation: %v", eventID, err)
			} else if promotedUser != nil {
				log.Printf("User %s (ID: %d) was promoted from waitlist for event %d.", promotedUser.Email, promotedUser.Id, eventID)
				// Potentially send notification here if not handled by ProcessNextOnWaitlist internally
			} else {
				log.Printf("No user promoted from waitlist for event %d, or waitlist was empty.", eventID)
			}
		}()
	} else {
		log.Printf("Event %d was not at full capacity. No need to process waitlist.", eventID)
	}

	return nil
}

func (s *eventService) GetRegisteredEvents(ctx context.Context, userID uuid.UUID) ([]model.Event, error) {
	return s.eventRepository.GetRegisteredEventByUserId(ctx, userID)
}

func (s *eventService) GetEventsByCategory(ctx context.Context, category string) ([]model.Event, error) {
	return s.eventRepository.GetEventsByCategory(ctx, category)
}

func (s *eventService) GetEventsByCriteria(ctx context.Context, keyword string, startDate string, endDate string) ([]model.Event, error) {
	return s.eventRepository.GetEventsByCriteria(ctx, keyword, startDate, endDate)
}
