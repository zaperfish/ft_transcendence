package event

import (
	// Std
	"context"
	"errors"
	"time"

	// Intern
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/errs"
	"ft_transcendence/backend/user"

	// Extern
	"github.com/danielgtaylor/huma/v2"
)

type EventHandler struct {
	service EventService
}

type EventDTO struct {
	ID              uint          `json:"id" doc:"ID of the event"`
	CreatedAt       time.Time     `json:"created_at" doc:"Time the event got created"`
	UpdatedAt       time.Time     `json:"updated_at" doc:"Time the event got updated"`
	Title           string        `json:"title" doc:"Name of the event"`
	Description     string        `json:"description" doc:"Description of the event"`
	StartTime       time.Time     `json:"start_time" doc:"Start time of the event"`
	Duration        int           `json:"duration" doc:"Duration of the event in minutes"`
	LocationName    string        `json:"location_name" doc:"Name of the location"`
	LocationAddress string        `json:"location_address" doc:"Address of the location"`
	MaxCapacity     uint          `json:"max_capacity" doc:"Maximum number of people the event supports"`
	NumRegistered   uint          `json:"num_registered" doc:"Number of people who registered for this event"`
	HasImage		bool		  `json:"has_image" doc:"Denotes if event has a custom image or not"`
	Self            *EventSelfDTO `json:"self,omitempty" doc:"Information about the authenticated user if authenticated"`
}

type EventSelfDTO struct {
	IsParticipant bool   `json:"is_participant" doc:"Shows if authenticated user is a participant of the event"`
	Role          string `json:"role" doc:"Shows user's role in the event"`
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
		HasImage:		 e.ImagePath != "",
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
		LocationAddress string    `json:"location_address" minLength:"5"  maxLength:"200" example:"Prinzessinnenstraße 19, 10969 Berlin" doc:"Address of the location"`
		MaxCapacity     uint      `json:"max_capacity"     minimum:"1"    maximum:"10000" example:"100"                                 doc:"Maximum number of attendees"`
	}
}

type CreateEventOutput struct {
	Body EventDTO
}

func (h *EventHandler) CreateEvent(ctx context.Context, input *CreateEventInput) (*CreateEventOutput, error) {
	userID, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("no authenticated user", err)
	}

	event := Event{
		Title:           input.Body.Title,
		Description:     input.Body.Description,
		StartTime:       input.Body.StartTime,
		Duration:        input.Body.Duration,
		LocationName:    input.Body.LocationName,
		LocationAddress: input.Body.LocationAddress,
		MaxCapacity:     input.Body.MaxCapacity,
	}

	created, err := h.service.CreateEventWithAdmin(ctx, &event, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to create event", err)
	}

	return &CreateEventOutput{Body: created.ToDTO()}, nil
}

type UpdateEventInput struct {
	ID   uint `path:"id" doc:"Event ID"`
	Body struct {
		Title           *string    `json:"title,omitempty"            minLength:"3"  maxLength:"100" example:"Go Meetup Berlin"                    doc:"Title of the event"`
		Description     *string    `json:"description,omitempty"      minLength:"10" maxLength:"500" example:"A monthly meetup for Go developers"  doc:"Description of the event"`
		StartTime       *time.Time `json:"start_time,omitempty"                                      example:"2026-06-15T18:00:00Z"                doc:"Start time of the event"`
		Duration        *int       `json:"duration,omitempty"         minimum:"15"   maximum:"480"   example:"120"                                 doc:"Duration of the event in minutes"`
		LocationName    *string    `json:"location_name,omitempty"    minLength:"3"  maxLength:"100" example:"Betahaus"                            doc:"Name of the location"`
		LocationAddress *string    `json:"location_address,omitempty" minLength:"5"  maxLength:"200" example:"Prinzessinnenstraße 19, 10969 Berlin" doc:"Address of the location"`
		MaxCapacity     *uint      `json:"max_capacity,omitempty"     minimum:"1"    maximum:"10000" example:"100"                                 doc:"Maximum number of attendees"`
	}
}

type UpdateEventOutput struct {
	Body EventDTO
}

func (h *EventHandler) UpdateEvent(ctx context.Context, input *UpdateEventInput) (*UpdateEventOutput, error) {
	if err := confirmAdminPriviliges(ctx, h, input.ID); err != nil {
		return nil, err
	}

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
	ID uint `path:"id" doc:"Event ID"`
}

type DeleteEventOutput struct {
	Body struct {
	}
}

func (h *EventHandler) DeleteEvent(ctx context.Context, input *DeleteEventInput) (*DeleteEventOutput, error) {
	if err := confirmAdminPriviliges(ctx, h, input.ID); err != nil {
		return nil, err
	}

	err := h.service.DeleteEvent(ctx, input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to delete event", err)
	}

	return &DeleteEventOutput{}, nil
}

type GetEventInput struct {
	ID uint `path:"id" doc:"Event ID"`
}

type GetEventOutput struct {
	Body EventDTO
}

func (h *EventHandler) GetEvent(ctx context.Context, input *GetEventInput) (*GetEventOutput, error) {
	userID, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("no authenticated user", err)
	}

	event, err := h.service.GetEventForUser(ctx, userID, input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to get event", err)
	}

	output := GetEventOutput{Body: event.ToDTO()}
	output.Body.Self = &EventSelfDTO{
		IsParticipant: event.IsParticipant,
		Role:          event.Role,
	}

	return &output, nil
}

type EventFilter string

const (
	EventFilterAll    EventFilter = "all"
	EventFilterMember EventFilter = "member"
	EventFilterAdmin  EventFilter = "admin"
)

type ListEventsInput struct {
	Page     int         `query:"page" minimum:"1" default:"1" doc:"Filter by page"`
	PageSize int         `query:"page_size" minimum:"1" default:"10" doc:"Page size"`
	Filter   EventFilter `query:"filter" enum:"all,member,admin" default:"all" doc:"Event filer"`
}

type ListEventsOutput struct {
	Body ListEventsOutputBody
}

type ListEventsOutputBody struct {
	Data     []EventDTO `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int64      `json:"total"`
}

func (h *EventHandler) ListEvents(ctx context.Context, input *ListEventsInput) (*ListEventsOutput, error) {
	userID, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("no authenticated user", err)
	}

	events, total, err := h.service.ListEvents(ctx, userID, input.PageSize, input.PageSize*(input.Page-1), input.Filter)
	if err != nil {
		return nil, huma.Error500InternalServerError("handler: failed to list events", err)
	}

	data := make([]EventDTO, len(events))
	for i, event := range events {
		data[i] = event.ToDTO()
		data[i].Self = &EventSelfDTO{
			IsParticipant: event.IsParticipant,
			Role:          event.Role,
		}
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
	EventID uint `path:"id" doc:"Event ID"`
	Body    struct {
		UserID uint   `json:"user_id" example:"1" doc:"ID of user to be added"`
		Role   string `json:"role" example:"member" doc:"member/admin"`
	}
}

type AddParticipantOutput struct {
	Body struct {
	}
}

func (h *EventHandler) AddParticipant(ctx context.Context, input *AddParticipantInput) (*AddParticipantOutput, error) {

	err := h.service.AddParticipantAs(ctx, input.EventID, input.Body.UserID, input.Body.Role)
	if err != nil {
		return nil, huma.Error500InternalServerError("", err)
	}

	return &AddParticipantOutput{}, nil
}

type RemoveParticipantInput struct {
	EventID uint `path:"eventID" doc:"Event ID"`
	UserID  uint `path:"userID" doc:"User ID"`
}

type RemoveParticipantOutput struct {
	Body struct {
	}
}

func (h *EventHandler) RemoveParticipant(ctx context.Context, input *RemoveParticipantInput) (*RemoveParticipantOutput, error) {
	if err := confirmAdminPriviliges(ctx, h, input.EventID); err != nil {
		return nil, err
	}

	err := h.service.RemoveParticipant(ctx, input.EventID, input.UserID)
	if err != nil && errors.Is(err, errs.ErrCanNotRemoveAdmin) {
		return nil, huma.Error403Forbidden(err.Error())
	}
	if err != nil && errors.Is(err, errs.ErrUserNotInEvent) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &RemoveParticipantOutput{}, nil
}

type ListParticipantsInput struct {
	EventID uint `path:"id" doc:"Event ID"`
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
		data[i] = *user.ToSummaryDTO()
	}

	return &ListParticipantsOutput{
		Body: ListParticipantsOutputBody{
			Data: data,
		},
	}, nil
}

// images
func (h *EventHandler) CreateImage(ctx context.Context, input *CreateImageInput) (*struct{}, error) {
	if err := confirmAdminPriviliges(ctx, h, input.EventID); err != nil {
		return nil, err
	}

	if err := h.service.CreateEventImage(ctx, input.EventID, input.Body, input.ContentType); err != nil {
		return nil, err
	}

	return nil, nil
}

type CreateImageInput struct {
	EventID uint `path:"id" doc:"Event ID"`
	ContentType string `header:"Content-Type"`
	Body []byte
}

func (h *EventHandler) GetImage(ctx context.Context, input *GetImageInput) (*GetImageOutput, error) {
	if err := confirmAdminPriviliges(ctx, h, input.EventID); err != nil {
		return nil, err
	}

	img, mime, err := h.service.GetEventImage(ctx, input.EventID)
	if err != nil {
		return nil, err
	}

	return &GetImageOutput{ContentType: mime, Body: img}, nil
}

type GetImageInput struct {
	EventID uint `path:"id" doc:"Event ID"`
}

type GetImageOutput struct {
	ContentType string `header:"Content-Type"`
	Body []byte
}

func (h *EventHandler) UpdateImage(ctx context.Context, input *UpdateImageInput) (*struct{}, error) {
	if err := confirmAdminPriviliges(ctx, h, input.EventID); err != nil {
		return nil, err
	}

	if err := h.service.CreateEventImage(ctx, input.EventID, input.Body, input.ContentType); err != nil {
		return nil, err
	}

	return nil, nil
}

type UpdateImageInput struct {
	EventID uint `path:"id" doc:"Event ID"`
	ContentType string `header:"Content-Type"`
	Body []byte
}

func (h *EventHandler) DeleteImage(ctx context.Context, input *DeleteImageInput) (*struct{}, error) {
	if err := confirmAdminPriviliges(ctx, h, input.EventID); err != nil {
		return nil, err
	}

	if err := h.service.DeleteEventImage(ctx, input.EventID); err != nil {
		return nil, err
	}

	return nil, nil
}

type DeleteImageInput struct {
	EventID uint `path:"id" doc:"Event ID"`
}

// helper
func confirmAdminPriviliges(ctx context.Context, h *EventHandler, eventID uint) error {
	userID, err := auth.UidFromCtx(ctx)
	if err != nil {
		return huma.Error401Unauthorized("no authenticated user", err)
	}
	event, err := h.service.GetEventForUser(ctx, userID, eventID)
	if err != nil && errors.Is(err, errs.ErrInternal) {
		return huma.Error500InternalServerError(err.Error())
	}
	if err != nil || event.Role != "admin" {
		return huma.Error401Unauthorized("must be admin")
	}
	return nil
}
