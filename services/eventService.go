package services

import (
	"context"
	"errors"
	"go-rest-api/model"
	"go-rest-api/repository"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	GetAllEvents(ctx context.Context) ([]model.Event, error)
	GetEventByID(ctx context.Context, id int64) (*model.Event, error)
	UpdateEvent(ctx context.Context, event *model.Event, userID int64) error
	DeleteEvent(ctx context.Context, id int64, userID int64) error
	RegisterForEvent(ctx context.Context, eventID, userID int64) error
	CancelEventRegistration(ctx context.Context, eventID, userID int64) error
}

type eventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepository repository.EventRepository) EventService {
	return &eventService{
		eventRepository: eventRepository,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, event *model.Event) error {
	return s.eventRepository.Save(ctx, event)
}

func (s *eventService) GetAllEvents(ctx context.Context) ([]model.Event, error) {
	return s.eventRepository.GetAllEvents(ctx)
}

func (s *eventService) GetEventByID(ctx context.Context, id int64) (*model.Event, error) {
	return s.eventRepository.GetEventById(ctx, id)
}

func (s *eventService) UpdateEvent(ctx context.Context, event *model.Event, userID int64) error {
	existingEvent, err := s.eventRepository.GetEventById(ctx, event.Id)
	if err != nil {
		return err
	}

	if existingEvent.UserIds != userID {
		return errors.New("unauthorized: you don't have permission to update this event")
	}

	return s.eventRepository.Update(ctx, event)
}

func (s *eventService) DeleteEvent(ctx context.Context, id int64, userID int64) error {
	existingEvent, err := s.eventRepository.GetEventById(ctx, id)
	if err != nil {
		return err
	}

	if existingEvent.UserIds != userID {
		return errors.New("unauthorized: you don't have permission to delete this event")
	}

	return s.eventRepository.DeleteEvent(ctx, id)
}

func (s *eventService) RegisterForEvent(ctx context.Context, eventID, userID int64) error {
	// Check if event exists
	_, err := s.eventRepository.GetEventById(ctx, eventID)
	if err != nil {
		return err
	}

	return s.eventRepository.RegisterEvent(ctx, eventID, userID)
}

func (s *eventService) CancelEventRegistration(ctx context.Context, eventID, userID int64) error {
	// Check if event exists
	_, err := s.eventRepository.GetEventById(ctx, eventID)
	if err != nil {
		return err
	}

	return s.eventRepository.CancelRegistration(ctx, eventID, userID)
}
