package event

import (
	"context"
	"errors"
	"fmt"

	"ft_transcendence/backend/user"

	"gorm.io/gorm"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	UpdateEvent(ctx context.Context, id uint, updates map[string]any) (*Event, error)
	DeleteEvent(ctx context.Context, id uint) error
	GetEvent(ctx context.Context, id uint) (*Event, error)
	ListEvents(ctx context.Context, limit, offset int) ([]Event, int64, error)
	ListEventsByUserID(ctx context.Context, limit, offset int, id uint) ([]Event, int64, error)
	AddParticipant(ctx context.Context, eventID, userID uint) error
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

func (s *eventServiceImpl) UpdateEvent(ctx context.Context, id uint, updates map[string]any) (*Event, error) {

	updated, err := s.repo.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *eventServiceImpl) DeleteEvent(ctx context.Context, id uint) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *eventServiceImpl) GetEvent(ctx context.Context, id uint) (*Event, error) {

	event, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventServiceImpl) ListEvents(ctx context.Context, limit, offset int) ([]Event, int64, error) {
	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	events, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (s *eventServiceImpl) ListEventsByUserID(ctx context.Context, limit, offset int, id uint) ([]Event, int64, error) {
	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	events, total, err := s.repo.ListByUserID(ctx, limit, offset, id)
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (s *eventServiceImpl) AddParticipant(ctx context.Context, eventID, userID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateParticipant(ctx, tx, eventID, userID); err != nil {
			return fmt.Errorf("failed to create participant: %w", err)
		}

		if err := s.repo.IncrementParticipantCount(ctx, tx, eventID, 1); err != nil {
			return fmt.Errorf("failed to increment participant count: %w", err)
		}

		return nil
	})
}

func (s *eventServiceImpl) RemoveParticipant(ctx context.Context, eventID, userID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DeleteParticipant(ctx, tx, eventID, userID); err != nil {
			return fmt.Errorf("failed to create participant: %w", err)
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
