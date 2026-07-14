package event

import (
	// Std
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"unsafe"

	// Internal
	"ft_transcendence/backend/errs"
	"ft_transcendence/backend/user"

	// External
	"gorm.io/gorm"
	"github.com/gabriel-vasile/mimetype"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	CreateEventWithAdmin(ctx context.Context, event *Event, userID uint) (*Event, error)
	UpdateEvent(ctx context.Context, userID uint, updates map[string]any) (*Event, error)
	DeleteEvent(ctx context.Context, userID uint) error
	GetEvent(ctx context.Context, eventID uint) (*Event, error)
	GetEventForUser(ctx context.Context, userID, eventID uint) (*EventWithUserContext, error)
	ListEvents(ctx context.Context, userID uint, limit, offset int, filer EventFilter) ([]EventWithUserContext, int64, error)
	AddParticipantAs(ctx context.Context, eventID, userID uint, role string) error
	RemoveParticipant(ctx context.Context, eventID, userID uint) error
	ListParticipants(ctx context.Context, eventID uint) ([]user.User, error)
	CreateEventImage(ctx context.Context, eventID uint, image []byte, contentType string) error
	GetEventImage(ctx context.Context, eventID uint) ([]byte, string, error)
	UpdateEventImage(ctx context.Context, eventID uint, image []byte, contentType string) error
	DeleteEventImage(ctx context.Context, eventID uint) error
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
	Role          string
}

func (s *eventServiceImpl) CreateEvent(ctx context.Context, e *Event) (*Event, error) {
	if len(e.Title) < 3 {
		return nil, errors.New("title must be at least 3 characters")
	}

	if e.Duration <= 0 {
		return nil, errors.New("duration must be greater than 0")
	}

	created, err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *eventServiceImpl) CreateEventWithAdmin(ctx context.Context, e *Event, userID uint) (*Event, error) {

	if len(e.Title) < 3 {
		return nil, errs.NewCamaError(errs.ErrInvalidInput, "title must be at least 3 characters")
	}

	if e.Duration <= 0 {
		return nil, errs.NewCamaError(errs.ErrInvalidInput, "duration must be greater than 0")
	}

	var created Event

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		userS := user.NewUserService(user.NewUserRepository(tx))
		if _, err := userS.GetUserByID(ctx, userID); err != nil {
			return err
		}

		ev, err := s.repo.Create(ctx, e)
		if err != nil {
			return err
		}

		err = s.repo.CreateParticipantAs(ctx, tx, ev.ID, userID, "admin")
		if err != nil {
			return err
		}

		created = *ev

		return nil

	}); err != nil {
		return nil, errs.ErrorDB(err)
	}

	return &created, nil
}

func (s *eventServiceImpl) UpdateEvent(ctx context.Context, eventID uint, updates map[string]any) (*Event, error) {

	updated, err := s.repo.Update(ctx, eventID, updates)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *eventServiceImpl) DeleteEvent(ctx context.Context, eventID uint) error {
	if err := s.DeleteEventImage(ctx, eventID); err != nil {
		log.Printf("DeleteEventImage: %v\n", err)
	}
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
		Event:         event.Event,
		IsParticipant: event.Role != "none",
		Role:          event.Role,
	}

	return &output, nil
}

func (s *eventServiceImpl) ListEvents(ctx context.Context, userID uint, limit, offset int, filter EventFilter) ([]EventWithUserContext, int64, error) {
	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	events, total, err := s.repo.ListByUserID(ctx, limit, offset, userID, filter)
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
		cap, err := s.repo.GetCapacity(ctx, tx, eventID)
		if err != nil {
			return err
		}

		count, err := s.repo.GetParticipantCount(ctx, tx, eventID)
		if err != nil {
			return err
		}

		if count >= cap {
			return errors.New("event full")
		}

		if err := s.repo.CreateParticipantAs(ctx, tx, eventID, userID, role); err != nil {
			return fmt.Errorf("failed to create participant: %w", err)
		}

		return nil
	})
}

func isValidRole(role string) bool {
	return role == "admin" || role == "member"
}

func (s *eventServiceImpl) RemoveParticipant(ctx context.Context, eventID, userID uint) error {
	event, err := s.repo.GetForUser(ctx, userID, eventID)
	if err != nil {
		return errors.New("failed to get event user information")
	}
	if event.Role == "none" {
		return errs.ErrUserNotInEvent
	}
	if event.Role == "admin" {
		return errs.ErrCanNotRemoveAdmin
	}

	if err := s.repo.DeleteParticipant(ctx, nil, eventID, userID); err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}
	return nil
}

func (s *eventServiceImpl) ListParticipants(ctx context.Context, eventID uint) ([]user.User, error) {
	users, err := s.repo.GetParticipants(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	return users, nil
}

var imagePathPrefix string = "/var/lib/ft_transcendence/images"
const maxImageSize = 104857600

// This should only work, when no image is associated with the event yet. This is why we can not use s.repo.Update() here as it would overwrite an existing image path. s.repo.CreateImagePath() makes sure to not overwrite an existing path.
func (s *eventServiceImpl) CreateEventImage(ctx context.Context, eventID uint, image []byte, contentType string) error {

	mtype := mimetype.Detect(image)

	if err := validateImage(image, contentType, mtype); err != nil {
		return errs.NewCamaError(errs.ErrInvalidInput, err.Error())
	}

	path := imagePathPrefix + "/" + strconv.FormatUint(uint64(eventID), 10) + mtype.Extension()
	if err := s.repo.CreateImagePath(ctx, eventID, path); err != nil {
		return errs.NewCamaError(errs.ErrInternal, err.Error())
	}

	if err := os.WriteFile(path, image, 0600); err != nil {
		s.repo.DeleteImagePath(ctx, eventID)
		return errs.NewCamaError(errs.ErrInternal, err.Error())
	}

	return nil
}

func (s *eventServiceImpl) GetEventImage(ctx context.Context, eventID uint) ([]byte, string, error) {
	path, err := s.repo.GetImagePath(ctx, eventID)
	if err != nil {
		return nil, "", errs.NewCamaError(errs.ErrNotFound, err.Error())
	}

	image, err := os.ReadFile(path)
	if err != nil {
		return nil, "", errs.NewCamaError(errs.ErrInternal, err.Error())
	}

	mtype := mimetype.Detect(image)

	return image, mtype.String(), nil
}

func (s *eventServiceImpl) UpdateEventImage(ctx context.Context, eventID uint, image []byte, contentType string) error {
	mtype := mimetype.Detect(image)

	if err := validateImage(image, contentType, mtype); err != nil {
		return errs.NewCamaError(errs.ErrInvalidInput, err.Error())
	}

	path, err := s.repo.GetImagePath(ctx, eventID)
	if err != nil {
		return errs.NewCamaError(errs.ErrNotFound, err.Error())
	}

	// note: WriteFile() atomically replaces the file at path
	err = os.WriteFile(path, image, 0600)
	if err != nil {
		return errs.NewCamaError(errs.ErrInternal, err.Error())
	}

	return nil
}

// when db update succeeds and image deletion fails there will 
func (s *eventServiceImpl) DeleteEventImage(ctx context.Context, eventID uint) error {
	path, err := s.repo.GetImagePath(ctx, eventID)
	if err != nil {
		return errs.NewCamaError(errs.ErrInvalidInput, err.Error())
	}

	err = s.repo.DeleteImagePath(ctx, eventID)
	if err != nil {
		return errs.NewCamaError(errs.ErrInternal, err.Error())
	}

	err = os.Remove(path)
	if err != nil {
		return errs.NewCamaError(errs.ErrInternal, err.Error())
	}

	return nil
}

func validateImage(image []byte, contentType string, mtype *mimetype.MIME) error {
	if unsafe.Sizeof(image) > maxImageSize {
		return errors.New("file too large")
	}

	if !mtype.Is(contentType) {
		return errors.New("Content-type header does not match file type")
	}

	if contentType != "image/png" {
		return errors.New("must be image/png")
	}

	return nil
}
