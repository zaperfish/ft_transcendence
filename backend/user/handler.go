package user

import (
    // Std
	"context"
	"errors"

    // Internal
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/errs"

    // External
	"github.com/danielgtaylor/huma/v2"
)

type UserHandler struct {
	s UserService
}

// register

func (h *UserHandler) handleRegisterUser(ctx context.Context, in *CreateInput) (*UserOutput, error) {
	u, err := h.s.CreateUser(ctx, in.Body)
	if errors.Is(err, errs.ErrInvalidInput) {
		return nil, huma.Error400BadRequest(err.Error())
	}
	if errors.Is(err, errs.ErrConflict) {
        return nil, huma.Error409Conflict(err.Error())
    }
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
    return &UserOutput{Body: *u}, nil
}

// login

func (h *UserHandler) handleLoginUser(ctx context.Context, in *LoginUserInput) (*LoginUserOutput, error) {
    u, err := h.getUserByName(ctx, in.Body.Name)
    if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error401Unauthorized(err.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
	
	match, err := auth.MatchPassword(in.Body.Password, u.PasswordHash)
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
    }
	if !match {
        return nil, huma.Error401Unauthorized(errs.ErrNotFound.Error())
	}

	cookie, err := auth.MakeJWTCookieFromID(u.ID)
    if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
    }

    out := &LoginUserOutput {
		SetCookie: cookie,
        Body: 	   u.ToSummaryDTO(),
    }

    return out, nil
}

// logout

func (h *UserHandler) handleLogoutUser(ctx context.Context, in *struct{}) (*LogoutUserOutput, error) {

    out := &LogoutUserOutput {
		SetCookie: auth.MakeLogoutCookie(),
    }

    return out, nil
}

// get

func (h *UserHandler) handleGetUser(ctx context.Context, in *GetUserInput) (*UserOutput, error) {
	u, err := h.s.GetUser(ctx, in.ID)
    if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error404NotFound(err.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
    return &UserOutput{Body: *u}, nil
}

// get list

func (h *UserHandler) handleGetUsers(ctx context.Context, in *GetUsersInput) (*UsersOutput, error) {

	us, err := h.s.GetUsers(ctx, in.Page, in.PageSize)
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
    
    out := UsersOutput {
        Body: UserListSummaryDTO {
            Data:       us,
            Page:       in.Page,
            PageSize:   in.PageSize,
            Total:      len(us),
        },
    }
	return &out, nil
}

// patch

func (h *UserHandler) handlePatchUser(ctx context.Context, in *PatchUserInput) (*UserOutput, error) {
	u, err := h.s.PatchUser(ctx, in.ID, in.Body)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &UserOutput{Body: *u}, nil
}

// patch password

func (h *UserHandler) handlePatchPassword(ctx context.Context, in *PatchPasswordInput) (*UserOutput, error) {
	u, err := h.s.PatchPassword(ctx, in.ID, in.Body)
	if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error404NotFound(err.Error())
	}
	if errors.Is(err, errs.ErrConflict) || errors.Is(err, errs.ErrInvalidInput) {
        return nil, huma.Error400BadRequest(err.Error())
	}
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UserOutput{Body: *u}, nil
}

// delete

func (h *UserHandler) handleDeleteUser(ctx context.Context, in *DeleteUserInput) (*struct{}, error) {
	err := h.s.DeleteUser(ctx, in.ID)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
    return nil, nil
}
