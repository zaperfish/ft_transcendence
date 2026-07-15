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
	"gorm.io/gorm"
)

type UserHandler struct {
	s UserService
}

func NewHandler(db *gorm.DB) UserHandler {
	return UserHandler{s: NewUserService(NewUserRepository(db))}
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
	u, cookie, err := h.s.LoginUser(ctx, in.Body.Name, in.Body.Password)
    if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error401Unauthorized(err.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
    out := &LoginUserOutput {
		SetCookie: cookie,
        Body: 	   *u,
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
	u, err := h.s.GetUserByID(ctx, in.ID)
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
