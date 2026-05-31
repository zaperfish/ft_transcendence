package event

import (
	"context"
	"errors"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	UpdateEvent(ctx context.Context, id string, updates map[string]any) (*Event, error)
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (*Event, error)
	ListEvents(ctx context.Context, limit, offset int) ([]*Event, error)
}

type eventServiceImpl struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) EventService {
	return &eventServiceImpl{repo: repo}
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

func (s *eventServiceImpl) UpdateEvent(ctx context.Context, id string, updates map[string]any) (*Event, error) {
	if id == "" {
		return nil, errors.New("missing event ID")
	}

	updated, err := s.repo.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *eventServiceImpl) DeleteEvent(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("missing event ID")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *eventServiceImpl) GetEvent(ctx context.Context, id string) (*Event, error) {
	if id == "" {
		return nil, errors.New("missing event ID")
	}

	event, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventServiceImpl) ListEvents(ctx context.Context, limit, offset int) ([]*Event, error) {
	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	events, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return events, nil
}
