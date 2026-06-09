package event

import (
	// Std
	"context"
	"errors"
	"fmt"

	// Internal
	"ft_transcendence/backend/user"

	// External
	"gorm.io/gorm"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	UpdateEvent(ctx context.Context, eventID uint, updates map[string]any) (*Event, error)
	DeleteEvent(ctx context.Context, eventID uint) error
	GetEvent(ctx context.Context, eventID uint) (*Event, error)
	ListEvents(ctx context.Context, userID uint, limit, offset int) ([]EventWithUserContext, int64, error)
	ListEventsByUserID(ctx context.Context, limit, offset int, userID uint) ([]Event, int64, error)
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
	Role string
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

func (s *eventServiceImpl) UpdateEvent(ctx context.Context, eventID uint, updates map[string]any) (*Event, error) {

	updated, err := s.repo.Update(ctx, eventID, updates)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// TODO: let DeleteParticipants be handled by GORM
func (s *eventServiceImpl) DeleteEvent(ctx context.Context, eventID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := s.repo.DeleteParticipants(ctx, eventID); err != nil {
			return err
		}

		if err := s.repo.Delete(ctx, eventID); err != nil {
			return err
		}

		return nil
	})
}

func (s *eventServiceImpl) GetEvent(ctx context.Context, eventID uint) (*Event, error) {

	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventServiceImpl) ListEvents(ctx context.Context, userID uint, limit, offset int) ([]EventWithUserContext, int64, error) {
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

	// participantEventIDs, err := s.repo.GetParticipantEventIDs(ctx, user_id)
	// if err != nil {
	// 	return nil, 0, err
	// }
	//
	// participantMap := make(map[uint]bool, len(participantEventIDs))
	// for _, id := range participantEventIDs {
	// 	participantMap[id] = true
	// }

	eventsWithUserCtx := make([]EventWithUserContext, 0, len(events))
	var role string
	for _, e := range events {

		role, err = s.repo.GetEventUsersRole(ctx, e.ID, userID) // save because strconv.IntSize
		if err != nil {
			return nil, 0, err
		}
		eventsWithUserCtx = append(eventsWithUserCtx, EventWithUserContext{
			Event:  e,
			Role: 	role,
		})
	}

	return eventsWithUserCtx, total, nil
}

func (s *eventServiceImpl) ListEventsByUserID(ctx context.Context, limit, offset int, userID uint) ([]Event, int64, error) {
	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	events, total, err := s.repo.ListByUserID(ctx, limit, offset, userID)
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (s *eventServiceImpl) AddParticipantAs(ctx context.Context, eventID, userID uint, role string) error {
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
