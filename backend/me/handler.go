package me

import (
    // Std
	"context"
	"errors"
	"net/http"

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

func NewHandler(db *gorm.DB, participantDisconnector event.ParticipantDisconnector) MeHandler {
	return MeHandler{su: user.NewUserService(user.NewUserRepository(db)),
					 se: event.NewEventService(event.NewEventRepository(db), db, participantDisconnector)}
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

type PatchMeDTO struct {
    Email *string    `json:"email,omitempty" example:"max@email.com" doc:"email address"`
}

type PatchMeInput struct {
	Body PatchMeDTO
}

func (h *MeHandler) handlePatchMe(ctx context.Context, in *PatchMeInput) (*user.UserOutput, error) {
	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	patch := user.PatchUserDTO{
		Email: in.Body.Email,
	}
	u, err := h.su.PatchUser(ctx, id, patch)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &user.UserOutput{Body: *u}, nil
}

// patch password me

type PatchMePasswordInput struct {
	Body user.PatchPasswordDTO
}

func (h *MeHandler) handlePatchPasswordMe(ctx context.Context, in *PatchMePasswordInput) (*user.UserOutput, error) {

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

type DeleteMeOutput struct {
	SetCookie http.Cookie 		`header:"Set-Cookie"`
}

func (h *MeHandler) handleDeleteMe(ctx context.Context, in *struct{}) (*DeleteMeOutput, error) {
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

    out := &DeleteMeOutput {
		SetCookie: auth.MakeLogoutCookie(),
    }

    return out, nil
}

// join event

type JoinEventInput struct {
	EventID uint `path:"id" doc:"Event ID"`
}

func (h *MeHandler) handleJoinEventMe(ctx context.Context, input *JoinEventInput) (*event.AddParticipantOutput, error) {
	uid, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	err = h.se.AddParticipantAs(ctx, input.EventID, uid, "member")
	if errors.Is(err, errs.ErrConflict) {
		return nil, huma.Error409Conflict(err.Error())
	}
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if errors.Is(err, errs.ErrInternal) || err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &event.AddParticipantOutput{}, nil
}

// leave event

type LeaveEventInput struct {
	EventID uint `path:"id" doc:"Event ID"`
}

func (h *MeHandler) handleLeaveEventMe(ctx context.Context, input *LeaveEventInput) (*event.AddParticipantOutput, error) {
	userID, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	err = h.se.RemoveParticipant(ctx, input.EventID, userID)
	if err != nil && errors.Is(err, errs.ErrCanNotRemoveAdmin) {
		return nil, huma.Error403Forbidden(err.Error())
	}
	if err != nil && errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}


	return &event.AddParticipantOutput{}, nil
}
