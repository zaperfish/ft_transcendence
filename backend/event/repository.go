package event

import (
	// Std
	"context"
	"fmt"
	"strconv"
	"time"

	// Intern
	"ft_transcendence/backend/errs"
	"ft_transcendence/backend/user"

	// Extern
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	Update(ctx context.Context, id string, updated map[string]any) (*Event, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*Event, error)
	List(ctx context.Context, limit, offset int) ([]Event, int64, error)
	ListByUserID(ctx context.Context, limit, offset int, id uint) ([]Event, int64, error)
	CreateParticipant(ctx context.Context, tx *gorm.DB, eventID, userID string) error
	DeleteParticipant(ctx context.Context, tx *gorm.DB, eventID, userID string) error
	IncrementParticipantCount(ctx context.Context, tx *gorm.DB, eventID string, amount int) error
	DecrementParticipantCount(ctx context.Context, tx *gorm.DB, eventID string, amount int) error
	GetParticipants(ctx context.Context, eventID string) ([]user.User, error)
	IsParticipant(ctx context.Context, eventID, userID string) (bool, error)
	GetParticipantEventIDs(ctx context.Context, userID string) ([]uint, error)
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

	Participants []user.User `gorm:"many2many:event_participants;joinForeignKey:EventID;joinReferences:UserID"`
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

func (r *eventRepositoryImpl) List(ctx context.Context, limit, offset int) ([]Event, int64, error) {
	models, err := gorm.G[GormEventModel](r.db.Debug()).Limit(limit).Offset(offset).Find(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve list of events: %w", err)
	}

	num_retrieved := len(models)
	events := make([]Event, num_retrieved)
	for i, model := range models {
		events[i] = *model.ToDomain()
	}

	var total int64
	gorm.G[GormEventModel](r.db.Debug()).
		Select("count(*)").
		Scan(ctx, &total)

	return events, total, nil
}

func (r *eventRepositoryImpl) ListByUserID(ctx context.Context, limit, offset int, id uint) ([]Event, int64, error) {

	var events []Event
	var count int64

	err := r.db.WithContext(ctx).
		Joins("JOIN event_participants ep ON ep.event_id = events.id").
		Where("ep.user_id = ?", id).
		Limit(limit).
		Offset(offset).
		Find(&events).
		Count(&count).Error

	if err != nil {
		return nil, 0, errs.ErrorDB(err)
	}

	return events, count, nil
}

func (r *eventRepositoryImpl) CreateParticipant(ctx context.Context, tx *gorm.DB, eventID, userID string) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	var count int64
	db.Table("event_participants").Where("event_id = ? AND user_id = ?", eventID, userID).Count(&count)
	if count > 0 {
		return fmt.Errorf("user is already participant")
	}

	event, err := gorm.G[GormEventModel](db.Debug()).Where("id = ?", eventID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to find event: %w", err)
	}

	user, err := gorm.G[user.User](db.Debug()).Where("id = ?", userID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	err = db.Model(&event).Association("Participants").Append(&user)
	if err != nil {
		return fmt.Errorf("failed to add user to event: %w", err)
	}

	return nil
}

func (r *eventRepositoryImpl) DeleteParticipant(ctx context.Context, tx *gorm.DB, eventID, userID string) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	var count int64
	db.Table("event_participants").Where("event_id = ? AND user_id = ?", eventID, userID).Count(&count)
	if count == 0 {
		return fmt.Errorf("user is not a participant")
	}

	event, err := gorm.G[GormEventModel](db.Debug()).Where("id = ?", eventID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to find event: %w", err)
	}

	user, err := gorm.G[user.User](db.Debug()).Where("id = ?", userID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	err = db.Model(&event).Association("Participants").Delete(&user)
	if err != nil {
		return fmt.Errorf("failed to add user to event: %w", err)
	}

	return nil
}

func (r *eventRepositoryImpl) IncrementParticipantCount(ctx context.Context, tx *gorm.DB, eventID string, amount int) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	rows, err := gorm.G[GormEventModel](db.Debug()).
		Where("id = ?", eventID).
		Where("num_registered + ? <= max_capacity", amount).
		Update(ctx, "num_registered", gorm.Expr("num_registered + ?", amount))
	if err != nil {
		return fmt.Errorf("failed to find event: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no rows updated, event maybe full")
	}

	return nil
}

func (r *eventRepositoryImpl) DecrementParticipantCount(ctx context.Context, tx *gorm.DB, eventID string, amount int) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	rows, err := gorm.G[GormEventModel](db.Debug()).
		Where("id = ?", eventID).
		Where("num_registered - ? >= 0", amount).
		Update(ctx, "num_registered", gorm.Expr("num_registered - ?", amount))
	if err != nil {
		return fmt.Errorf("failed to find event: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (r *eventRepositoryImpl) GetParticipants(ctx context.Context, eventID string) ([]user.User, error) {
	var models []user.User

	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN event_participants ep ON ep.user_id = users.id").
		Where("ep.event_id = ?", eventID).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	return models, nil
}

func (r *eventRepositoryImpl) IsParticipant(ctx context.Context, eventID, userID string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Table("event_participants").
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *eventRepositoryImpl) GetParticipantEventIDs(ctx context.Context, userID string) ([]uint, error) {
	var participantEventIDs []uint
	if userID == " " {

		fmt.Println("HEEEEEEENOTOGODD")
	}

	fmt.Println("g")
	fmt.Println(userID, "A")
	fmt.Println("g")

	err := r.db.WithContext(ctx).
		Table("event_participants").
		Where("user_id = ?", userID).
		Pluck("event_id", &participantEventIDs).Error

	if err != nil {
		return nil, err
	}

	return participantEventIDs, nil
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
