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

	// Associations
	Labels []Label `gorm:"many2many:event_labels;"`
}

type PostEventDTO struct {
	Title           string    `json:"title"            minLength:"3"  maxLength:"100" example:"Go Meetup Berlin"                    doc:"Title of the event"`
	Description     string    `json:"description"      minLength:"10" maxLength:"500" example:"A monthly meetup for Go developers"  doc:"Description of the event"`
	StartTime       time.Time `json:"start_time"                                      example:"2026-06-15T18:00:00Z"                doc:"Start time of the event"`
	Duration        int       `json:"duration"         minimum:"15"   maximum:"480"   example:"120"                                 doc:"Duration of the event in minutes"`
	LocationName    string    `json:"location_name"    minLength:"3"  maxLength:"100" example:"Betahaus"                            doc:"Name of the location"`
	LocationAddress string    `json:"location_address" minLength:"5"  maxLength:"200" example:"Prinzessinnenstraße 19, 10969 Berlin"doc:"Address of the location"`
	MaxCapacity     int       `json:"max_capacity"     minimum:"1"    maximum:"10000" example:"100"                                 doc:"Maximum number of attendees"`
}

type EventInput struct {
	Body PostEventDTO
}

type EventDTO struct {
	ID              uint       `json:"id" doc:"ID of the event"`
	CreatedAt       time.Time  `json:"created_at" doc:"Time the event got created"`
	UpdatedAt       time.Time  `json:"updated_at" doc:"Time the event got updated"`
	Title           string     `json:"title" doc:"Name of the event"`
	Description     string     `json:"description" doc:"Description of the event"`
	StartTime       time.Time  `json:"start_time" doc:"Start time of the event"`
	Duration        int        `json:"duration" doc:"Duration of the event in minutes"`
	LocationName    string     `json:"location_name" doc:"Name of the location"`
	LocationAddress string     `json:"location_address" doc:"Address of the location"`
	MaxCapacity     int        `json:"max_capacity" doc:"Maximum number of people the event supports"`
	NumRegistered   int        `json:"num_registered" doc:"Number of people who registered for this event"`
	Labels          []LabelDTO `json:"labels" doc:"Labels of the event"`
}

type PatchEventInput struct {
	ID   string `path:"id" doc:"Event ID"`
	Body PatchEventDTO
}
type PatchEventDTO struct {
	Title           *string    `json:"title,omitempty"            minLength:"3"  maxLength:"100" example:"Go Meetup Berlin"                 doc:"Title of the event"`
	Description     *string    `json:"description,omitempty"      minLength:"10" maxLength:"500" example:"A monthly meetup for Go developers" doc:"Description of the event"`
	StartTime       *time.Time `json:"start_time,omitempty"                                      example:"2026-06-15T18:00:00Z"               doc:"Start time of the event"`
	Duration        *int       `json:"duration,omitempty"         minimum:"15"   maximum:"480"   example:"120"                                doc:"Duration in minutes"`
	LocationName    *string    `json:"location_name,omitempty"    minLength:"3"  maxLength:"100" example:"Betahaus"                           doc:"Name of the location"`
	LocationAddress *string    `json:"location_address,omitempty" minLength:"5"  maxLength:"200" example:"Prinzessinnenstraße 19, 10969 Berlin" doc:"Address of the location"`
	MaxCapacity     *int       `json:"max_capacity,omitempty"     minimum:"1"    maximum:"10000" example:"100"                                doc:"Maximum number of attendees"`
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

func eventToEventDTO(event *Event) EventDTO {
	labels := make([]LabelDTO, len(event.Labels))

	for i, label := range event.Labels {
		labels[i] = labelToLabelDTO(&label)
	}

	eventDTO := EventDTO{
		ID:              event.ID,
		CreatedAt:       event.CreatedAt,
		UpdatedAt:       event.UpdatedAt,
		Title:           event.Title,
		Description:     event.Description,
		StartTime:       event.StartTime,
		Duration:        event.Duration,
		LocationName:    event.LocationName,
		LocationAddress: event.LocationAddress,
		MaxCapacity:     event.MaxCapacity,
		NumRegistered:   event.NumRegistered,
		Labels:          labels,
	}

	return eventDTO
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

type PathID struct {
	ID string `path:"id" doc:"ID"`
}

func preloadAll(db gorm.PreloadBuilder) error {
	return nil
}

func (h *EventHandler) HandleGetEventByID(ctx context.Context, input *PathID) (*EventOutput, error) {
	event, err := gorm.G[Event](h.db.Debug()).Preload("Labels", nil).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	eventOutput := &EventOutput{
		Body: eventToEventDTO(&event),
	}

	return eventOutput, nil
}

func (h *EventHandler) HandlePatchEvent(ctx context.Context, input *PatchEventInput) (*EventOutput, error) {
	updates := map[string]any{}

	if input.Body.Title != nil {
		updates["title"] = *input.Body.Title
	}
	if input.Body.Description != nil {
		updates["description"] = *input.Body.Description
	}
	if input.Body.StartTime != nil {
		updates["start_time"] = *input.Body.StartTime
	}
	if input.Body.Duration != nil {
		updates["duration"] = *input.Body.Duration
	}
	if input.Body.LocationName != nil {
		updates["location_name"] = *input.Body.LocationName
	}
	if input.Body.LocationAddress != nil {
		updates["location_address"] = *input.Body.LocationAddress
	}
	if input.Body.MaxCapacity != nil {
		updates["max_capacity"] = *input.Body.MaxCapacity
	}

	_, err := gorm.G[map[string]interface{}](h.db.Debug()).Table("events").Where("id = ?", input.ID).Updates(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to save patched event: %w", err)
	}

	updated, err := gorm.G[Event](h.db.Debug()).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated event: %w", err)
	}

	return &EventOutput{Body: eventToEventDTO(&updated)}, nil

}

type GetEventsInput struct {
	Page     int `query:"page" minimum:"1" default:"1" doc:"Filter by page"`
	PageSize int `query:"page_size" minimum:"1" default:"10" doc:"Page size"`
}

func (h *EventHandler) HandleGetEvents(ctx context.Context, input *GetEventsInput) (*EventsOutput, error) {
	base := gorm.G[Event](h.db.Debug())

	offset := (input.Page - 1) * input.PageSize
	q := base.Preload("Labels", nil)
	q = q.Limit(input.PageSize)
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

func (h *EventHandler) HandleDeleteEvent(ctx context.Context, input *PathID) (*EmptyResponse, error) {
	rows_affected, err := gorm.G[Event](h.db.Debug()).Where("id = ?", input.ID).Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete event: %w", err)
	}

	if rows_affected == 0 {
		return nil, fmt.Errorf("failed to delete event: record not found")
	}

	return &EmptyResponse{}, nil
}

type AddLabelInput struct {
	EventID string `path:"id" doc:"ID"`
	Body    struct {
		LabelID int `json:"label_id" doc:"id of the label to add"`
	}
}

func (h *EventHandler) HandleAddLabel(ctx context.Context, input *AddLabelInput) (*EmptyResponse, error) {
	event, err := gorm.G[Event](h.db.Debug()).Where("id = ?", input.EventID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find event: %w", err)
	}

	label, err := gorm.G[Label](h.db.Debug()).Where("id = ?", input.Body.LabelID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find label: %w", err)
	}

	err = h.db.Model(&event).Association("Labels").Append(&label)
	if err != nil {
		return nil, fmt.Errorf("failed to append label to event: %w", err)
	}

	return nil, nil
}

type DeleteLabelInput struct {
	EventID string `path:eventID`
	LabelID string `path:labelID`
}

func (h *EventHandler) HandleDeleteLabel(ctx context.Context, input *DeleteLabelInput) (*EmptyResponse, error) {
	event, err := gorm.G[Event](h.db.Debug()).Where("id = ?", input.EventID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find event: %w", err)
	}

	label, err := gorm.G[Label](h.db.Debug()).Where("id = ?", input.LabelID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find label: %w", err)
	}

	err = h.db.Model(&event).Association("Labels").Delete(&label)
	if err != nil {
		return nil, fmt.Errorf("failed to append label to event: %w", err)
	}

	return &EmptyResponse{}, nil
}

type LabelHandler struct {
	db *gorm.DB
}

type Label struct {
	gorm.Model

	Name string `gorm:"not null;uniqueIndex;check:length(name) >= 2"`
}

type LabelInput struct {
	Body CreateLabelDTO
}

type LabelOutput struct {
	Body LabelDTO
}

type CreateLabelDTO struct {
	Name string `json:"name" minLenght:"2" maxLength:"15" example:"Go" doc:"Name of the label"`
}

type LabelDTO struct {
	ID        uint      `json:"id" doc:"ID of the event"`
	CreatedAt time.Time `json:"created_at" doc:"Time the event got created"`
	UpdatedAt time.Time `json:"updated_at" doc:"Time the event got updated"`
	Name      string    `json:"name" doc:"Name of the label"`
}

func labelToLabelDTO(label *Label) LabelDTO {
	labelDTO := LabelDTO{
		Name:      label.Name,
		ID:        label.ID,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}

	return labelDTO
}

func (h *LabelHandler) HandlePostLabel(ctx context.Context, input *LabelInput) (*LabelOutput, error) {
	label := Label{
		Name: input.Body.Name,
	}

	err := gorm.G[Label](h.db.Debug()).Create(ctx, &label)
	if err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}

	labelOutput := LabelOutput{
		Body: labelToLabelDTO(&label),
	}

	return &labelOutput, nil
}

func (h *LabelHandler) HandleGetLabelByID(ctx context.Context, input *PathID) (*LabelOutput, error) {
	label, err := gorm.G[Label](h.db.Debug()).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}

	labelOutput := LabelOutput{
		Body: labelToLabelDTO(&label),
	}

	return &labelOutput, nil
}

type LabelListDTO struct {
	Data     []LabelDTO `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
}

type LabelsOutput struct {
	Body LabelListDTO
}

type GetLabelsInput struct {
	Page     int `query:"page" minimum:"1" default:"1" doc:"Filter by page"`
	PageSize int `query:"page_size" minimum:"1" default:"10" doc:"Page size"`
}

func (h *EventHandler) HandleGetLabels(ctx context.Context, input *GetLabelsInput) (*LabelsOutput, error) {
	base := gorm.G[Label](h.db.Debug())

	offset := (input.Page - 1) * input.PageSize
	q := base.Limit(input.PageSize)
	q = q.Offset(offset)

	g_labels, err := q.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	total := len(g_labels)
	labels := make([]LabelDTO, total)
	for i, g_label := range g_labels {
		labels[i] = labelToLabelDTO(&g_label)
	}

	labelsOutput := &LabelsOutput{
		Body: LabelListDTO{
			Data:     labels,
			Page:     input.Page,
			PageSize: input.PageSize,
			Total:    total,
		},
	}

	return labelsOutput, nil
}

func (h *LabelHandler) HandleDeleteLabel(ctx context.Context, input *PathID) (*EmptyResponse, error) {
	rows_affected, err := gorm.G[Label](h.db.Debug()).Where("id = ?", input.ID).Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete label: %w", err)
	}

	if rows_affected == 0 {
		return nil, fmt.Errorf("failed to delete label: record not found")
	}

	return &EmptyResponse{}, nil
}

func RegisterApi(api huma.API, db *gorm.DB) {
	event_handler := EventHandler{db: db}
	label_handler := LabelHandler{db: db}

	db.AutoMigrate(&Event{})
	db.AutoMigrate(&Label{})

	// Register POST /events
	huma.Register(api, huma.Operation{
		OperationID:   "create-event",
		Method:        http.MethodPost,
		Path:          "/api/events",
		Summary:       "Create an event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusCreated,
	}, event_handler.HandlePostEvent)

	// Register GET /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "get-event-by-id",
		Method:        http.MethodGet,
		Path:          "/api/events/{id}",
		Summary:       "Get an event by ID",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, event_handler.HandleGetEventByID)

	// Register GET /events
	huma.Register(api, huma.Operation{
		OperationID:   "get-events",
		Method:        http.MethodGet,
		Path:          "/api/events",
		Summary:       "Get a list of events",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, event_handler.HandleGetEvents)

	// Register DELETE /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "delete-event",
		Method:        http.MethodDelete,
		Path:          "/api/events/{id}",
		Summary:       "Delete an event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, event_handler.HandleDeleteEvent)

	// Register PATCH /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "patch-event",
		Method:        http.MethodPatch,
		Path:          "/api/events/{id}",
		Summary:       "Patch an event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, event_handler.HandlePatchEvent)

	// Register POST /events/{id}/labels
	huma.Register(api, huma.Operation{
		OperationID:   "add-label-to-event",
		Method:        http.MethodPost,
		Path:          "/api/events/{id}/labels",
		Summary:       "Add a label to an event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusNoContent,
	}, event_handler.HandleAddLabel)

	// Register DELETE /events/{eventID}/labels/{labelID}
	huma.Register(api, huma.Operation{
		OperationID:   "delete-label-from-event",
		Method:        http.MethodDelete,
		Path:          "/api/events/{eventID}/labels/{labelID}",
		Summary:       "Delete a label from an event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusNoContent,
	}, event_handler.HandleDeleteLabel)

	// Register POST /labels
	huma.Register(api, huma.Operation{
		OperationID:   "create-label",
		Method:        http.MethodPost,
		Path:          "/api/labels",
		Summary:       "Create a label",
		Tags:          []string{"Labels"},
		DefaultStatus: http.StatusCreated,
	}, label_handler.HandlePostLabel)

	// Register GET /labels/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "get-label",
		Method:        http.MethodGet,
		Path:          "/api/labels/{id}",
		Summary:       "Get a label by ID",
		Tags:          []string{"Labels"},
		DefaultStatus: http.StatusOK,
	}, label_handler.HandleGetLabelByID)

	// Register GET /labels
	huma.Register(api, huma.Operation{
		OperationID:   "get-labels",
		Method:        http.MethodGet,
		Path:          "/api/labels",
		Summary:       "Get a list of labels",
		Tags:          []string{"Labels"},
		DefaultStatus: http.StatusOK,
	}, event_handler.HandleGetLabels)

	// Register DELETE /labels/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "delete-label",
		Method:        http.MethodDelete,
		Path:          "/api/labels/{id}",
		Summary:       "Delete a a label",
		Tags:          []string{"Labels"},
		DefaultStatus: http.StatusOK,
	}, event_handler.HandleDeleteLabel)

}
