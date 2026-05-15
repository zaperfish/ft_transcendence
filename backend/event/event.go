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

type EventHandler struct {
	db *gorm.DB
}

type Event struct {
	gorm.Model
	Name          string
	MaxCapacity   int
	NumRegistered int
}

type EventInput struct {
	Body struct {
		Name          string `json:"name" doc:"Name of the event"`
		MaxCapacity   int    `json:"max_capacity" doc:"Maximum number of people the event supports"`
		NumRegistered int    `json:"num_registered" doc:"Number of people who registered for this event"`
	}
}

type EventOutput struct {
	Body EventBody
}

type EventBody struct {
	ID            uint      `json:"id" doc:"ID of the event"`
	CreatedAt     time.Time `json:"created_at" doc:"Time the event got created"`
	UpdatedAt     time.Time `json:"updated_at" doc:"Time the event got updated"`
	Name          string    `json:"name" doc:"Name of the event"`
	MaxCapacity   int       `json:"max_capacity" doc:"Maximum number of people the event supports"`
	NumRegistered int       `json:"num_registered" doc:"Number of people who registered for this event"`
}

func (h *EventHandler) HandlePostEvent(ctx context.Context, input *EventInput) (*EventOutput, error) {
	event := Event{
		Name:          input.Body.Name,
		MaxCapacity:   input.Body.MaxCapacity,
		NumRegistered: input.Body.NumRegistered,
	}

	err := gorm.G[Event](h.db).Create(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	eventOutput := &EventOutput{
		Body: EventBody{
			ID:            event.ID,
			CreatedAt:     event.CreatedAt,
			UpdatedAt:     event.UpdatedAt,
			Name:          event.Name,
			MaxCapacity:   event.MaxCapacity,
			NumRegistered: event.NumRegistered,
		},
	}

	return eventOutput, nil
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

	// huma.Register(api, huma.Operation{
	// 	OperationID:   "get-event",
	// 	Method:        http.MethodPost,
	// 	Path:          "/api/events/{id}",
	// 	Summary:       "Get an event",
	// 	Tags:          []string{"Events"},
	// 	DefaultStatus: http.StatusCreated,
	// }, h.HandlePostEvent)

	// huma.Register(api, huma.Operation{
	// 	OperationID:   "get-event",
	// 	Method:        http.MethodPost,
	// 	Path:          "/api/events/{id}",
	// 	Summary:       "Get an event",
	// 	Tags:          []string{"Events"},
	// 	DefaultStatus: http.StatusCreated,
	// }, h.HandlePostEvent)
}
