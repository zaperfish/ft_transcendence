package event

import (
	// Std
	"context"
	"fmt"
	"strconv"
	"time"

	// Intern
	"ft_transcendence/backend/errs"

	// Extern
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	Update(ctx context.Context, id string, updated map[string]any) (*Event, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*Event, error)
	List(ctx context.Context, limit, offset int) ([]*Event, error)
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepositoryImpl{db: db}
}

type GormEventModel struct {
	gorm.Model

	Title           string    `gorm:"type:varchar(255);not null"`
	Description     string    `gorm:"type:text"`
	StartTime       time.Time `gorm:"not null"`
	Duration        int       `gorm:"type:smallint;not null"`
	LocationName    string    `gorm:"type:varchar(255)"`
	LocationAddress string    `gorm:"type:varchar(255)"`
	MaxCapacity     int       `gorm:"not null;"`
	NumRegistered   int       `gorm:"not null;"`
}

func (GormEventModel) TableName() string {
	return "events"
}

func (m *GormEventModel) ToDomain() *Event {
	return &Event{
		ID:              strconv.FormatUint(uint64(m.ID), 10),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		Title:           m.Title,
		Description:     m.Description,
		StartTime:       m.StartTime,
		Duration:        m.Duration,
		LocationName:    m.LocationName,
		LocationAddress: m.LocationAddress,
		MaxCapacity:     m.MaxCapacity,
		NumRegistered:   m.NumRegistered,
	}
}

func (r *eventRepositoryImpl) Create(ctx context.Context, event *Event) (*Event, error) {
	model := GormEventModel{
		Title:           event.Title,
		Description:     event.Description,
		StartTime:       event.StartTime,
		Duration:        event.Duration,
		LocationName:    event.LocationName,
		LocationAddress: event.LocationAddress,
		MaxCapacity:     event.MaxCapacity,
		NumRegistered:   event.NumRegistered,
	}

	err := gorm.G[GormEventModel](r.db.Debug()).Create(ctx, &model)
	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *eventRepositoryImpl) Update(ctx context.Context, id string, updates map[string]any) (*Event, error) {
	rows, err := gorm.G[map[string]any](r.db.Debug()).Table("events").Where("id = ?", id).Updates(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	if rows == 0 {
		return nil, fmt.Errorf("no event updated")
	}

	model, err := gorm.G[GormEventModel](r.db.Debug()).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated event: %w", err)
	}

	return model.ToDomain(), nil
}

func (r *eventRepositoryImpl) Delete(ctx context.Context, id string) error {
	rows, err := gorm.G[GormEventModel](r.db.Debug()).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no event deleted")
	}

	return nil
}

func (r *eventRepositoryImpl) Get(ctx context.Context, id string) (*Event, error) {
	model, err := gorm.G[GormEventModel](r.db.Debug()).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve event: %w", err)
	}

	return model.ToDomain(), nil
}

func (r *eventRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*Event, error) {
	models, err := gorm.G[GormEventModel](r.db.Debug()).Limit(limit).Offset(offset).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve list of events: %w", err)
	}

	total := len(models)
	events := make([]*Event, total)
	for i, model := range models {
		events[i] = model.ToDomain()
	}

	return events, nil
}

func IsParticipant(ctx context.Context, db *gorm.DB, eventID uint, userID uint) (bool, error) {
	var count int64

	err := db.WithContext(ctx).
		Table("event_participants").
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count).Error

	if err != nil {
		return false, errs.ErrorDB(err)
	}

	return count > 0, nil
}

func EventExists(ctx context.Context, db *gorm.DB, eventID uint) (bool, error) {
	var count int64

	err := db.WithContext(ctx).
		Model(&Event{}).
		Where("id = ?", eventID).
		Count(&count).Error

	if err != nil {
		return false, errs.ErrorDB(err)
	}

	return count > 0, nil
}
