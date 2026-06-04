package event

import (
	// Std
	"context"
	"ft_transcendence/backend/user"
	"strconv"
	"time"

	// Intern
	// "ft_transcendence/backend/user"

	// Extern
	"github.com/danielgtaylor/huma/v2"
)

type EventHandler struct {
	service EventService
}

type EventDTO struct {
	ID              string    `json:"id" doc:"ID of the event"`
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

func (e *Event) ToDTO() EventDTO {
	eventDTO := EventDTO{
		ID:              e.ID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		Title:           e.Title,
		Description:     e.Description,
		StartTime:       e.StartTime,
		Duration:        e.Duration,
		LocationName:    e.LocationName,
		LocationAddress: e.LocationAddress,
		MaxCapacity:     e.MaxCapacity,
		NumRegistered:   e.NumRegistered,
	}

	return eventDTO
}

func NewEventHandler(service EventService) EventHandler {
	return EventHandler{service: service}
}

type CreateEventInput struct {
	Body struct {
		Title           string    `json:"title"            minLength:"3"  maxLength:"100" example:"Go Meetup Berlin"                    doc:"Title of the event"`
		Description     string    `json:"description"      minLength:"10" maxLength:"500" example:"A monthly meetup for Go developers"  doc:"Description of the event"`
		StartTime       time.Time `json:"start_time"                                      example:"2026-06-15T18:00:00Z"                doc:"Start time of the event"`
		Duration        int       `json:"duration"         minimum:"15"   maximum:"480"   example:"120"                                 doc:"Duration of the event in minutes"`
		LocationName    string    `json:"location_name"    minLength:"3"  maxLength:"100" example:"Betahaus"                            doc:"Name of the location"`
		LocationAddress string    `json:"location_address" minLength:"5"  maxLength:"200" example:"Prinzessinnenstraße 19, 10969 Berlin"doc:"Address of the location"`
		MaxCapacity     int       `json:"max_capacity"     minimum:"1"    maximum:"10000" example:"100"                                 doc:"Maximum number of attendees"`
	}
}

type CreateEventOutput struct {
	Body EventDTO
}

func (h *EventHandler) CreateEvent(ctx context.Context, input *CreateEventInput) (*CreateEventOutput, error) {
	event := Event{
		Title:           input.Body.Title,
		Description:     input.Body.Description,
		StartTime:       input.Body.StartTime,
		Duration:        input.Body.Duration,
		LocationName:    input.Body.LocationName,
		LocationAddress: input.Body.LocationAddress,
		MaxCapacity:     input.Body.MaxCapacity,
	}

	created, err := h.service.CreateEvent(ctx, &event)
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to create event", err)
	}

	return &CreateEventOutput{Body: created.ToDTO()}, nil
}

type UpdateEventInput struct {
	ID   string `path:"id" doc:"Event ID"`
	Body struct {
		Title           *string    `json:"title,omitempty"            minLength:"3"  maxLength:"100" example:"Go Meetup Berlin"                    doc:"Title of the event"`
		Description     *string    `json:"description,omitempty"      minLength:"10" maxLength:"500" example:"A monthly meetup for Go developers"  doc:"Description of the event"`
		StartTime       *time.Time `json:"start_time,omitempty"                                      example:"2026-06-15T18:00:00Z"                doc:"Start time of the event"`
		Duration        *int       `json:"duration,omitempty"         minimum:"15"   maximum:"480"   example:"120"                                 doc:"Duration of the event in minutes"`
		LocationName    *string    `json:"location_name,omitempty"    minLength:"3"  maxLength:"100" example:"Betahaus"                            doc:"Name of the location"`
		LocationAddress *string    `json:"location_address,omitempty" minLength:"5"  maxLength:"200" example:"Prinzessinnenstraße 19, 10969 Berlin"doc:"Address of the location"`
		MaxCapacity     *int       `json:"max_capacity,omitempty"     minimum:"1"    maximum:"10000" example:"100"                                 doc:"Maximum number of attendees"`
	}
}

type UpdateEventOutput struct {
	Body EventDTO
}

func (h *EventHandler) UpdateEvent(ctx context.Context, input *UpdateEventInput) (*UpdateEventOutput, error) {
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

	updated, err := h.service.UpdateEvent(ctx, input.ID, updates)
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to update event", err)
	}

	return &UpdateEventOutput{Body: updated.ToDTO()}, nil
}

type DeleteEventInput struct {
	ID string `path:"id" doc:"Event ID"`
}

type DeleteEventOutput struct {
	Body struct {
	}
}

func (h *EventHandler) DeleteEvent(ctx context.Context, input *DeleteEventInput) (*DeleteEventOutput, error) {
	err := h.service.DeleteEvent(ctx, input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to delete event", err)
	}

	return &DeleteEventOutput{}, nil
}

type GetEventInput struct {
	ID string `path:"id" doc:"Event ID"`
}

type GetEventOutput struct {
	Body EventDTO
}

func (h *EventHandler) GetEvent(ctx context.Context, input *GetEventInput) (*GetEventOutput, error) {
	event, err := h.service.GetEvent(ctx, input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to get event", err)
	}

	return &GetEventOutput{Body: event.ToDTO()}, nil
}

type ListEventsInput struct {
	Page     int `query:"page" minimum:"1" default:"1" doc:"Filter by page"`
	PageSize int `query:"page_size" minimum:"1" default:"10" doc:"Page size"`
}

type ListEventsOutput struct {
	Body ListEventsOutputBody
}

type ListEventsOutputBody struct {
	Data     []EventDTO `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
}

func (h *EventHandler) ListEvents(ctx context.Context, input *ListEventsInput) (*ListEventsOutput, error) {
	events, err := h.service.ListEvents(ctx, input.PageSize, input.PageSize*(input.Page-1))
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to list events", err)
	}

	total := len(events)
	data := make([]EventDTO, total)
	for i, event := range events {
		data[i] = event.ToDTO()
	}

	return &ListEventsOutput{
		Body: ListEventsOutputBody{
			Data:     data,
			Page:     input.Page,
			PageSize: input.PageSize,
			Total:    total,
		},
	}, nil
}

type AddParticipantInput struct {
	EventID string `path:"id" doc:"Event ID"`
	Body    struct {
		UserID int `json:"user_id"`
	}
}

type AddParticipantOutput struct {
	Body struct {
	}
}

func (h *EventHandler) AddParticipant(ctx context.Context, input *AddParticipantInput) (*AddParticipantOutput, error) {
	userID := strconv.Itoa(input.Body.UserID)

	err := h.service.AddParticipant(ctx, input.EventID, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("", err)
	}

	return &AddParticipantOutput{}, nil
}

type RemoveParticipantInput struct {
	EventID string `path:"eventID" doc:"Event ID"`
	UserID  string `path:"userID" doc:"User ID"`
}

type RemoveParticipantOutput struct {
	Body struct {
	}
}

func (h *EventHandler) RemoveParticipant(ctx context.Context, input *RemoveParticipantInput) (*RemoveParticipantOutput, error) {
	err := h.service.RemoveParticipant(ctx, input.EventID, input.UserID)
	if err != nil {
		return nil, huma.Error500InternalServerError("", err)
	}

	return &RemoveParticipantOutput{}, nil
}

type ListParticipantsInput struct {
	EventID string `path:"id" doc:"Event ID"`
}

type ListParticipantsOutput struct {
	Body ListParticipantsOutputBody
}

type ListParticipantsOutputBody struct {
	Data []user.UserSummaryDTO `json:"data"`
}

func (h *EventHandler) ListParticipants(ctx context.Context, input *ListParticipantsInput) (*ListParticipantsOutput, error) {
	users, err := h.service.ListParticipants(ctx, input.EventID)
	if err != nil {
		return nil, huma.Error500InternalServerError("", err)
	}

	total := len(users)
	data := make([]user.UserSummaryDTO, total)
	for i, user := range users {
		data[i] = user.ToSummaryDTO()
	}

	return &ListParticipantsOutput{
		Body: ListParticipantsOutputBody{
			Data: data,
		},
	}, nil
}
