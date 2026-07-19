package event

import (
	// Std
	"context"
	"time"
	"log"

	// Intern
	"ft_transcendence/backend/errs"
	"ft_transcendence/backend/eventusers"
	"ft_transcendence/backend/user"

	// Extern
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	Update(ctx context.Context, eventID uint, updated map[string]any) (*Event, error)
	Delete(ctx context.Context, eventID uint) error
	DeleteParticipants(ctx context.Context, eventID uint) error
	Get(ctx context.Context, eventID uint) (*Event, error)
	GetCapacity(ctx context.Context, tx *gorm.DB, eventID uint) (uint, error)
	GetForUser(ctx context.Context, userID, eventID uint) (*EventWithRole, error)
	GetParticipantCount(ctx context.Context, tx *gorm.DB, eventID uint) (uint, error)
	List(ctx context.Context, limit, offset int) ([]Event, int64, error)
	ListByUserID(ctx context.Context, limit, offset int, userID uint, filter EventFilter) ([]EventWithRole, int64, error)
	CreateParticipantAs(ctx context.Context, tx *gorm.DB, eventID, userID uint, role string) error
	DeleteParticipant(ctx context.Context, tx *gorm.DB, eventID, userID uint) error
	GetParticipants(ctx context.Context, eventID uint) ([]user.User, error)
	GetParticipantRole(ctx context.Context, eventID, userID uint) (bool, string, error)
	GetParticipantEventIDs(ctx context.Context, userID uint) ([]uint, error)
	CreateImagePath(ctx context.Context, eventID uint, path string) error
	GetImagePath(ctx context.Context, eventID uint) (string, error)
	DeleteImagePath(ctx context.Context, eventID uint) error
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
	MaxCapacity     uint      `gorm:"not null;"`
	ImagePath     	string    `gorm:"type:varchar(255)"`
}

func (GormEventModel) TableName() string {
	return "events"
}

func (m *GormEventModel) ToDomain() *Event {
	return &Event{
		ID:              m.ID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		Title:           m.Title,
		Description:     m.Description,
		StartTime:       m.StartTime,
		Duration:        m.Duration,
		LocationName:    m.LocationName,
		LocationAddress: m.LocationAddress,
		MaxCapacity:     m.MaxCapacity,
		ImagePath:		 m.ImagePath,
		NumRegistered:   0,
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
	}

	err := gorm.G[GormEventModel](r.db.Debug()).Create(ctx, &model)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	return model.ToDomain(), nil
}

// after delete hook
func (e *GormEventModel) AfterDelete(tx *gorm.DB) error {
	return tx.
		Model(&eventusers.EventUser{}).
		Where("event_id = ?", e.ID).
		Delete(&eventusers.EventUser{}).Error
}

func (r *eventRepositoryImpl) Update(ctx context.Context, eventID uint, updates map[string]any) (*Event, error) {

	count, err := r.GetParticipantCount(ctx, nil, eventID)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	_, err = gorm.G[map[string]any](r.db.Debug()).Table("events").Where("id = ?", eventID).Updates(ctx, updates)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	model, err := gorm.G[GormEventModel](r.db.Debug()).Where("id = ?", eventID).First(ctx)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	ret := model.ToDomain()
	ret.NumRegistered = count

	return model.ToDomain(), nil
}

func (r *eventRepositoryImpl) Delete(ctx context.Context, eventID uint) error {
	event := GormEventModel{
		Model: gorm.Model{ID: eventID},
	}
	result := r.db.Debug().Delete(&event)
	if result.Error != nil {
		log.Printf("db error: %v\n", result.Error)
		return errs.ErrorDB(result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("no rows affected\n")
		return errs.NewCamaError(errs.ErrNotFound, "")
	}

	return nil
}

func (r *eventRepositoryImpl) DeleteParticipants(ctx context.Context, eventID uint) error {
	err := r.db.
		WithContext(ctx).
		Where("event_id = ?", eventID).
		Delete(&eventusers.EventUser{}).
		Error
	return errs.ErrorDB(err)
}

func (r *eventRepositoryImpl) Get(ctx context.Context, eventID uint) (*Event, error) {
	model, err := gorm.G[GormEventModel](r.db.Debug()).Where("id = ?", eventID).First(ctx)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	count, err := r.GetParticipantCount(ctx, nil, eventID)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	ret := model.ToDomain()
	ret.NumRegistered = count

	return model.ToDomain(), nil
}

func (r *eventRepositoryImpl) GetForUser(ctx context.Context, userID, eventID uint) (*EventWithRole, error) {
	count, err := r.GetParticipantCount(ctx, nil, eventID)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	var eventRole EventWithRole
	tx := r.db.WithContext(ctx).
		Model(&GormEventModel{}).
		Select(`
        events.*,
        COALESCE(event_users.role, 'none') as role
    `).
		Joins(`
        LEFT JOIN event_users 
        ON event_users.event_id = events.id
        AND event_users.user_id = ?
		AND event_users.deleted_at IS NULL
    `, userID).
		Where("events.id = ?", eventID).
		Scan(&eventRole)
	if tx.Error != nil {
		return nil, errs.ErrorDB(tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, errs.NewCamaError(errs.ErrNotFound, "")
	}

	eventRole.NumRegistered = count

	return &eventRole, nil
}

func (r *eventRepositoryImpl) GetCapacity(ctx context.Context, tx *gorm.DB, eventID uint) (uint, error) {
	if tx == nil {
		tx = r.db
	}
	var cap uint
	ret := tx.WithContext(ctx).
		Model(&Event{}).
		Select("max_capacity").
		Where("id = ?", eventID).
		Scan(&cap)
	if ret.Error != nil {
		return 0, errs.ErrorDB(ret.Error)
	}
	if ret.RowsAffected == 0 {
		return 0, errs.NewCamaError(errs.ErrNotFound, "")
	}
	return cap, nil
}

func (r *eventRepositoryImpl) GetParticipantCount(ctx context.Context, tx *gorm.DB, eventID uint) (uint, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).
		Model(&eventusers.EventUser{}).
		Where("event_id = ?", eventID).
		Count(&count).Error; err != nil {
		return 0, errs.ErrorDB(err)
	}

	if count < 0 || uint64(count) > uint64(^uint(0)) {
		return 0, errs.NewCamaError(errs.ErrInternal, "")
	}

	return uint(count), nil
}

// having to individually query the participant count for each event is the tradeof for getting rid of NumParticipants in the GormEventModel
// there might be a smarter, more sql heavy version to do this with a single query but I don't know how
func (r *eventRepositoryImpl) List(ctx context.Context, limit, offset int) ([]Event, int64, error) {
	models, err := gorm.G[GormEventModel](r.db.Debug()).Limit(limit).Offset(offset).Find(ctx)
	if err != nil {
		return nil, 0, errs.ErrorDB(err)
	}

	num_retrieved := len(models)
	events := make([]Event, num_retrieved)
	for i, model := range models {
		events[i] = *model.ToDomain()
		count, err := r.GetParticipantCount(ctx, nil, events[i].ID)
		if err != nil {
			return nil, 0, errs.ErrorDB(err)
		}
		events[i].NumRegistered = count
	}

	var total int64
	gorm.G[GormEventModel](r.db.Debug()).
		Select("count(*)").
		Scan(ctx, &total)

	return events, total, nil
}

type EventWithRole struct {
	Event
	Role string
}

// having to individually query the participant count for each event is the tradeof for getting rid of NumParticipants in the GormEventModel
// there might be a smarter, more sql heavy version to do this with a single query but I don't know how
func (r *eventRepositoryImpl) ListByUserID(ctx context.Context, limit, offset int, userID uint, filter EventFilter) ([]EventWithRole, int64, error) {
	if _, err := user.NewUserRepository(r.db).GetByID(ctx, userID); err != nil {
		return nil, 0, err
	}

	q := r.db.WithContext(ctx).
		Model(&GormEventModel{}).
		Select(`
        events.*,
        COALESCE(event_users.role, 'none') as role
    `).
		Where("events.start_time > ?", time.Now().UTC()).
		Joins(`
        LEFT JOIN event_users 
        ON event_users.event_id = events.id 
		AND event_users.user_id = ?
		AND event_users.deleted_at IS NULL
    `, userID)

	if q.Error != nil {
		return nil, 0, errs.ErrorDB(q.Error)
	}

	switch filter {
	case EventFilterMember:
		q = q.Where("event_users.role = ?", "member")

	case EventFilterAdmin:
		q = q.Where("event_users.role = ?", "admin")

	case EventFilterAll:
		// No filter
	}

	var eventsRoles []EventWithRole
	if err := q.
		Limit(limit).
		Offset(offset).
		Scan(&eventsRoles).Error; err != nil {
		return nil, 0, errs.ErrorDB(err)
	}

	for i := range eventsRoles {
		count, err := r.GetParticipantCount(ctx, nil, eventsRoles[i].ID)
		if err != nil {
			return nil, 0, errs.ErrorDB(err)
		}
		eventsRoles[i].NumRegistered = count
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&GormEventModel{}).
		Count(&count).Error; err != nil {
		return nil, 0, errs.ErrorDB(err)
	}

	return eventsRoles, count, nil
}

func (r *eventRepositoryImpl) CreateParticipantAs(ctx context.Context, tx *gorm.DB, eventID, userID uint, role string) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	event, err := gorm.G[GormEventModel](db.Debug()).Where("id = ?", eventID).First(ctx)
	if err != nil {
		return errs.ErrorDB(err)
	}
	if event.StartTime.Before(time.Now().UTC()) {
		return errs.NewCamaError(errs.ErrConflict, "event already expired")
	}

	_, err = gorm.G[user.User](db.Debug()).Where("id = ?", userID).First(ctx)
	if err != nil {
		return errs.ErrorDB(err)
	}

	var count int64
	err = db.Model(&eventusers.EventUser{}).Where("event_id = ? AND user_id = ?", eventID, userID).Count(&count).Error
	if err != nil {
		return errs.ErrorDB(err)
	}
	if count > 0 {
		return errs.NewCamaError(errs.ErrConflict, "user already subscribed to event")
	}

	err = db.WithContext(ctx).Create(&eventusers.EventUser{
		UserID:  userID,
		EventID: eventID,
		Role:    role,
	}).Error

	if err != nil {
		return errs.ErrorDB(err)
	}

	return nil
}

func (r *eventRepositoryImpl) DeleteParticipant(ctx context.Context, tx *gorm.DB, eventID, userID uint) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	var count int64
	err := db.Model(&eventusers.EventUser{}).
		WithContext(ctx).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count).
		Error
	if err != nil {
		return errs.ErrorDB(err)
	}
	if count <= 0 {
		return errs.NewCamaError(errs.ErrNotFound, "user is not a participant")
	}

	err = db.
		WithContext(ctx).
		Model(&eventusers.EventUser{}).
		Where("user_id = ?", userID).
		Where("event_id = ?", eventID).
		Delete(&eventusers.EventUser{}).
		Error
	if err != nil {
		return errs.ErrorDB(err)
	}

	return nil
}

// TODO: also get roles
func (r *eventRepositoryImpl) GetParticipants(ctx context.Context, eventID uint) ([]user.User, error) {
	var models []user.User

	err := r.db.WithContext(ctx).
		Model(&user.User{}).
		Joins(`
			JOIN event_users eu 
			ON eu.user_id = users.id
			AND eu.deleted_at IS NULL`).
		Where("eu.event_id = ?", eventID).
		Find(&models).Error

	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	return models, nil
}

func (r *eventRepositoryImpl) GetParticipantRole(ctx context.Context, eventID, userID uint) (bool, string, error) {
	var count int64
	var role eventusers.EventUser
	if err := r.db.WithContext(ctx).
		Model(&eventusers.EventUser{}).
		Where("user_id = ?", userID).
		Where("event_id = ?", eventID).
		Count(&count).
		Scan(&role).Error; err != nil || count == 0 {
		return false, "none", errs.ErrorDB(err)
	}

	return true, role.Role, nil
}

// TODO: rework or remove
func (r *eventRepositoryImpl) GetParticipantEventIDs(ctx context.Context, userID uint) ([]uint, error) {
	var participantEventIDs []uint

	err := r.db.WithContext(ctx).
		Model(&eventusers.EventUser{}).
		Where("user_id = ?", userID).
		Pluck("event_id", &participantEventIDs).Error

	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	return participantEventIDs, nil
}

func (r *eventRepositoryImpl) CreateImagePath(ctx context.Context, eventID uint, path string) error {

	_, err := gorm.G[GormEventModel](r.db.Debug()).Where("id = ?", eventID).Where("image_path IS NULL or image_path = ''").Update(ctx, "image_path", path)
	if err != nil {
		return errs.ErrorDB(err)
	}

	return nil
}

func (r *eventRepositoryImpl) GetImagePath(ctx context.Context, eventID uint) (string, error) {
	var path string
	err := r.db.Model(&GormEventModel{}).Select("image_path").Where("id = ?", eventID).Scan(&path).Error
	if err != nil {
		return "", errs.ErrorDB(err)
	}
	if path == "" {
		return "", errs.NewCamaError(errs.ErrNotFound, "")
	}
	return path, nil
}

func (r *eventRepositoryImpl) DeleteImagePath(ctx context.Context, eventID uint) error {
	_, err := gorm.G[GormEventModel](r.db.Debug()).Where("id = ?", eventID).Update(ctx, "image_path", "")
	if err != nil {
		return errs.ErrorDB(err)
	}

	return nil
}

func IsParticipant(ctx context.Context, db *gorm.DB, eventID uint, userID uint) (bool, error) {
	var count int64

	err := db.WithContext(ctx).
		Model(&eventusers.EventUser{}).
		Where("user_id = ?", userID).
		Where("event_id = ?", eventID).
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
