package event

import (
	// Std
	"context"
	"fmt"
	"net/http"
	"time"

	// Extern
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type EmptyResponse struct{}

type EventHandler struct {
	db *gorm.DB
}

type Event struct {
	gorm.Model

	// Core
	Title       string `gorm:"not null;check:length(title) >= 3"`
	Description string
	StartTime   time.Time
	Duration    int

	// Location
	LocationName    string
	LocationAddress string

	// Capacity
	MaxCapacity   int `gorm:"check:max_capacity >= 0"`
	NumRegistered int `gorm:"check:max_capacity >= 0"`
}

type CreateEventDTO struct {
	Title           string    `json:"title"            minLength:"3"  maxLength:"100" example:"Go Meetup Berlin"                    doc:"Title of the event"`
	Description     string    `json:"description"      minLength:"10" maxLength:"500" example:"A monthly meetup for Go developers"  doc:"Description of the event"`
	StartTime       time.Time `json:"start_time"                                      example:"2026-06-15T18:00:00Z"                doc:"Start time of the event"`
	Duration        int       `json:"duration"         minimum:"15"   maximum:"480"   example:"120"                                 doc:"Duration of the event in minutes"`
	LocationName    string    `json:"location_name"    minLength:"3"  maxLength:"100" example:"Betahaus"                            doc:"Name of the location"`
	LocationAddress string    `json:"location_address" minLength:"5"  maxLength:"200" example:"Prinzessinnenstraße 19, 10969 Berlin"doc:"Address of the location"`
	MaxCapacity     int       `json:"max_capacity"     minimum:"1"    maximum:"10000" example:"100"                                 doc:"Maximum number of attendees"`
}

type EventInput struct {
	Body CreateEventDTO
}

type EventDTO struct {
	ID              uint      `json:"id" doc:"ID of the event"`
	CreatedAt       time.Time `json:"created_at" doc:"Time the event got created"`
	UpdatedAt       time.Time `json:"updated_at" doc:"Time the event got updated"`
	Title           string    `json:"title" doc:"Name of the event"`
	Description     string    `json:"description" doc:"Description of the event"`
	StartTime       time.Time `json:"start_time" doc:"Start time of the event"`
	Duration        int       `json:"duration" doc:"Duration of the event in minutes"`
	LocationName    string    `json:"location_name" doc:"Name of the location"`
	LocationAddress string    `json:"location_address" doc:"Address of the location"`
	MaxCapacity     int       `json:"max_capacity" doc:"Maximum number of people the event supports"`
	NumRegistered   int       `json:"num_registered" doc:"Number of people who registered for this event"`
}

type EventOutput struct {
	Body EventDTO
}

type EventListDTO struct {
	Data     []EventDTO `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
}

type EventsOutput struct {
	Body EventListDTO
}

func eventToEventDTO(gorm_event *Event) EventDTO {
	event := EventDTO{
		ID:              gorm_event.ID,
		CreatedAt:       gorm_event.CreatedAt,
		UpdatedAt:       gorm_event.UpdatedAt,
		Title:           gorm_event.Title,
		Description:     gorm_event.Description,
		StartTime:       gorm_event.StartTime,
		Duration:        gorm_event.Duration,
		LocationName:    gorm_event.LocationName,
		LocationAddress: gorm_event.LocationAddress,
		MaxCapacity:     gorm_event.MaxCapacity,
		NumRegistered:   gorm_event.NumRegistered,
	}

	return event
}

func (h *EventHandler) HandlePostEvent(ctx context.Context, input *EventInput) (*EventOutput, error) {
	event := Event{
		Title:           input.Body.Title,
		Description:     input.Body.Description,
		StartTime:       input.Body.StartTime,
		Duration:        input.Body.Duration,
		LocationName:    input.Body.LocationName,
		LocationAddress: input.Body.LocationAddress,
		MaxCapacity:     input.Body.MaxCapacity,
		NumRegistered:   0,
	}

	err := gorm.G[Event](h.db).Create(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	eventOutput := &EventOutput{
		Body: eventToEventDTO(&event),
	}

	return eventOutput, nil
}

type EventID struct {
	ID string `path:"id" doc:"ID of the event"`
}

func (h *EventHandler) HandleGetEventByID(ctx context.Context, input *EventID) (*EventOutput, error) {
	event, err := gorm.G[Event](h.db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	eventOutput := &EventOutput{
		Body: eventToEventDTO(&event),
	}

	return eventOutput, nil
}

type GetEventsInput struct {
	Page     int `query:"page" minimum:"1" default:"1" doc:"Filter by page"`
	PageSize int `query:"page_size" minimum:"1" default:"10" doc:"Page size"`
}

func (h *EventHandler) HandleGetEvents(ctx context.Context, input *GetEventsInput) (*EventsOutput, error) {
	base := gorm.G[Event](h.db.Debug())

	offset := (input.Page - 1) * input.PageSize
	q := base.Limit(input.PageSize)
	q = q.Offset(offset)

	g_events, err := q.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	total := len(g_events)
	events := make([]EventDTO, total)
	for i, g_event := range g_events {
		events[i] = eventToEventDTO(&g_event)
	}

	eventsOutput := &EventsOutput{
		Body: EventListDTO{
			Data:     events,
			Page:     input.Page,
			PageSize: input.PageSize,
			Total:    total,
		},
	}

	return eventsOutput, nil
}

func (h *EventHandler) HandleDeleteEvent(ctx context.Context, input *EventID) (*EmptyResponse, error) {
	rows_affected, err := gorm.G[Event](h.db.Debug()).Where("id = ?", input.ID).Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete event: %w", err)
	}

	if rows_affected == 0 {
		return nil, fmt.Errorf("failed to delete event: record not found")
	}

	return &EmptyResponse{}, nil
}

func RegisterApi(api huma.API, db *gorm.DB) {
	h := EventHandler{db: db}

	db.AutoMigrate(&Event{})

	// Register POST /events
	huma.Register(api, huma.Operation{
		OperationID:   "post-event",
		Method:        http.MethodPost,
		Path:          "/api/events",
		Summary:       "Post an event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusCreated,
	}, h.HandlePostEvent)

	// Register GET /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "get-event-by-id",
		Method:        http.MethodGet,
		Path:          "/api/events/{id}",
		Summary:       "Get an event by ID",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetEventByID)

	// Register GET /events
	huma.Register(api, huma.Operation{
		OperationID:   "get-events",
		Method:        http.MethodGet,
		Path:          "/api/events",
		Summary:       "Get a list of events",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetEvents)

	// Register DELETE /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "delete-event",
		Method:        http.MethodDelete,
		Path:          "/api/events/{id}",
		Summary:       "Delete an event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, h.HandleDeleteEvent)

	// // Register POST /events/{id}/participants
	// huma.Register(api, huma.Operation{
	// 	OperationID:   "join-event",
	// 	Method:        http.MethodPost,
	// 	Path:          "/api/events/{id}/participants",
	// 	Summary:       "Join an event",
	// 	Tags:          []string{"Events"},
	// 	DefaultStatus: http.StatusOK,
	// }, h.HandleJoinEvent)
}
