package me

import (
    // Std
	"context"
	"errors"

    // Internal
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/errs"
	"ft_transcendence/backend/event"
	"ft_transcendence/backend/user"

    // External
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// get me

type MeHandler struct {
	su user.UserService
	se event.EventService
}

func NewHandler(db *gorm.DB) MeHandler {
	return MeHandler{su: user.NewUserService(user.NewUserRepository(db)),
					 se: event.NewEventService(event.NewEventRepository(db), db)}
}

func (h *MeHandler) handleGetMe(ctx context.Context, in *struct{}) (*user.UserOutput, error) {
	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	u, err := h.su.GetUserByID(ctx, id)
    if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error404NotFound(err.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
    return &user.UserOutput{Body: *u}, nil
}

// patch me

func (h *MeHandler) handlePatchMe(ctx context.Context, in *user.PatchUserInput) (*user.UserOutput, error) {
	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	u, err := h.su.PatchUser(ctx, id, in.Body)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &user.UserOutput{Body: *u}, nil
}

// patch password me

func (h *MeHandler) handlePatchPasswordMe(ctx context.Context, in *user.PatchPasswordInput) (*user.UserOutput, error) {

	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	u, err := h.su.PatchPassword(ctx, id, in.Body)
	if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error404NotFound(err.Error())
	}
	if errors.Is(err, errs.ErrConflict) || errors.Is(err, errs.ErrInvalidInput) {
        return nil, huma.Error400BadRequest(err.Error())
	}
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
	return &user.UserOutput{Body: *u}, nil
}

// delete me

func (h *MeHandler) handleDeleteMe(ctx context.Context, in *struct{}) (*struct{}, error) {
	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	err = h.su.DeleteUser(ctx, id)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
    return nil, nil
}

// join event

type JoinEventInput struct {
	EventID string `path:"id" doc:"Event ID"`
}

func (h *MeHandler) handleJoinEventMe(ctx context.Context, input *JoinEventInput) (*event.AddParticipantOutput, error) {
	sub, err := auth.ClaimFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	err = h.se.AddParticipant(ctx, input.EventID, sub)
	if err != nil {
		return nil, huma.Error500InternalServerError("", err)
	}

	return &event.AddParticipantOutput{}, nil
}

// create event

// func (h *MeHandler) handleCreateEventMe(ctx context.Context,  input *event.CreateEventInput) (*event.CreateEventOutput, error) {
// 	id, err := auth.UidFromCtx(ctx)
// 	if err != nil {
// 		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
// 	}
//
// 	event := event.Event{
// 		Title:           input.Body.Title,
// 		Description:     input.Body.Description,
// 		StartTime:       input.Body.StartTime,
// 		Duration:        input.Body.Duration,
// 		LocationName:    input.Body.LocationName,
// 		LocationAddress: input.Body.LocationAddress,
// 		MaxCapacity:     input.Body.MaxCapacity,
// 	}
//
// 	created, err := h.se.CreateEvent(ctx, &event)
// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("handler: failed to create event", err)
// 	}
//
// 	err := h.service.AddParticipant(ctx, input.EventID, userID)
// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("", err)
// 	}
//
// }

//
// func (h *Handler) handleEventsMe(ctx context.Context, in *struct{}) (*event.ListEventsOutput, error) {
// 	id, err := auth.UidFromCtx(ctx)
// 	if err != nil {
// 		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
// 	}
// }
//
// func (h *Handler) handleAdminEventsMe(ctx context.Context. in *struct{}) (*event.ListEventsOutput, error) {
// 	id, err := auth.UidFromCtx(ctx)
// 	if err != nil {
// 		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
// 	}
// }
