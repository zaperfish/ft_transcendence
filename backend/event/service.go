package event

import (
	"context"
	"errors"
	"fmt"

	"ft_transcendence/backend/errs"
	"ft_transcendence/backend/user"

	"gorm.io/gorm"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	CreateEventWithAdmin(ctx context.Context, event *Event, userID uint) (*Event, error)
	UpdateEvent(ctx context.Context, userID uint, updates map[string]any) (*Event, error)
	DeleteEvent(ctx context.Context, userID uint) error
	GetEvent(ctx context.Context, eventID uint) (*Event, error)
	GetEventForUser(ctx context.Context, userID, eventID uint) (*EventWithUserContext, error)
	ListEvents(ctx context.Context, userID uint, limit, offset int) ([]EventWithUserContext, int64, error)
	AddParticipantAs(ctx context.Context, eventID, userID uint, role string) error
	RemoveParticipant(ctx context.Context, eventID, userID uint) error
	ListParticipants(ctx context.Context, eventID uint) ([]user.User, error)
}

type eventServiceImpl struct {
	repo EventRepository
	db   *gorm.DB
}

func NewEventService(repo EventRepository, db *gorm.DB) EventService {
	return &eventServiceImpl{repo: repo, db: db}
}

type EventWithUserContext struct {
	Event
	IsParticipant bool
	Role		  string
}

func (s *eventServiceImpl) CreateEvent(ctx context.Context, e *Event) (*Event, error) {
	if len(e.Title) < 3 {
		return nil, errors.New("title must be at least 3 characters")
	}

	if e.Duration <= 0 {
		return nil, errors.New("duration must be greater than 0")
	}

	if e.MaxCapacity < 0 {
		return nil, errors.New("max capacity cannot be negative")
	}

	e.NumRegistered = 0

	created, err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *eventServiceImpl) CreateEventWithAdmin(ctx context.Context, e *Event, userID uint) (*Event, error) {

	if len(e.Title) < 3 {
		return nil, errors.New("title must be at least 3 characters")
	}

	if e.Duration <= 0 {
		return nil, errors.New("duration must be greater than 0")
	}

	if e.MaxCapacity < 0 {
		return nil, errors.New("max capacity cannot be negative")
	}

	e.NumRegistered = 0

	var created Event

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		userS := user.NewUserService(user.NewUserRepository(tx))
		if _, err := userS.GetUserByID(ctx, userID); err != nil {
			return err
		}

		created, err := s.repo.Create(ctx, e)
		if err != nil {
			return err
		}

		err = s.repo.CreateParticipantAs(ctx, tx, created.ID, userID, "admin")
		if err != nil {
			return err
		}

		if err := s.repo.IncrementParticipantCount(ctx, tx, created.ID, 1); err != nil {
			return fmt.Errorf("failed to increment participant count: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *eventServiceImpl) UpdateEvent(ctx context.Context, id uint, updates map[string]any) (*Event, error) {

	updated, err := s.repo.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *eventServiceImpl) DeleteEvent(ctx context.Context, eventID uint) error {
	return s.repo.Delete(ctx, eventID)
}

func (s *eventServiceImpl) GetEvent(ctx context.Context, eventID uint) (*Event, error) {

	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventServiceImpl) GetEventForUser(ctx context.Context, userID, eventID uint) (*EventWithUserContext, error) {

	event, err := s.repo.GetForUser(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}

	output := EventWithUserContext{
		Event: event.Event,
		IsParticipant: event.Role != "none",
		Role: event.Role,
	}

	return &output, nil
}

func (s *eventServiceImpl) ListEvents(ctx context.Context, userID uint, limit, offset int) ([]EventWithUserContext, int64, error) {
	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	events, total, err := s.repo.ListByUserID(ctx, limit, offset, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("service: failed to list: %w", err)
	}

	eventsWithUserCtx := make([]EventWithUserContext, len(events))
	for i, e := range events {
		eventsWithUserCtx[i].Event = e.Event
		eventsWithUserCtx[i].IsParticipant = e.Role != "none"
		eventsWithUserCtx[i].Role = e.Role
	}

	return eventsWithUserCtx, total, nil
}

func (s *eventServiceImpl) AddParticipantAs(ctx context.Context, eventID, userID uint, role string) error {
	if !isValidRole(role) {
		return errs.ErrInvalidInput
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateParticipantAs(ctx, tx, eventID, userID, role); err != nil {
			return fmt.Errorf("failed to create participant: %w", err)
		}

		if err := s.repo.IncrementParticipantCount(ctx, tx, eventID, 1); err != nil {
			return fmt.Errorf("failed to increment participant count: %w", err)
		}

		return nil
	})
}

func isValidRole(role string) bool {
	return role == "admin" || role == "member"
}

func (s *eventServiceImpl) RemoveParticipant(ctx context.Context, eventID, userID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DeleteParticipant(ctx, tx, eventID, userID); err != nil {
			return fmt.Errorf("failed to remove participant: %w", err)
		}

		if err := s.repo.DecrementParticipantCount(ctx, tx, eventID, 1); err != nil {
			return fmt.Errorf("failed to decrement participant count: %w", err)
		}

		return nil
	})
}

func (s *eventServiceImpl) ListParticipants(ctx context.Context, eventID uint) ([]user.User, error) {
	users, err := s.repo.GetParticipants(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	return users, nil
}
