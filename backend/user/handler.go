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

// login

func (h *Handler) handleLoginUser(ctx context.Context, in *loginUserInput) (*LoginUserOutput, error) {
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

// register

func (h *Handler) handleRegisterUser(ctx context.Context, in *createInput) (*userOutput, error) {

	if err := validateParameters(&in.Body); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	hash, err := auth.CreateHash(in.Body.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

    u := User {
        Name:       	in.Body.Name,
        Email:      	in.Body.Email,
        PasswordHash:   hash,
    }

    err = h.creatUser(ctx, &u)
    if errors.Is(err, errs.ErrConflict) {
        return nil, huma.Error409Conflict(err.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}

    return &userOutput{Body: u.ToSummaryDTO()}, nil
}

func validateParameters(u *CreateUserDTO) error {
	if err := auth.ValidUserName(u.Name); err != nil {
		return err
	}
	if err := auth.ValidUserEmail(u.Email); err != nil {
		return err
	}
	if err := auth.ValidUserPassword(u.Password); err != nil {
		return err
	}
	return nil
}

// logout

func (h *Handler) handleLogoutUser(ctx context.Context, in *struct{}) (*LogoutUserOutput, error) {

    out := &LogoutUserOutput {
		SetCookie: auth.MakeLogoutCookie(),
    }

    return out, nil
}

// get

func (h *Handler) handleGetUser(ctx context.Context, in *getUserInput) (*userOutput, error) {
	u, err := h.getUserByID(ctx, in.ID)
    if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error404NotFound(err.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}

    return &userOutput{Body: u.ToSummaryDTO()}, nil
}

// get list

func (h *Handler) handleGetUsers(ctx context.Context, in *getUsersInput) (*usersOutput, error) {

	us, err := h.getUsersList(ctx, UserFilter(*in))
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
    
    userList := make([]UserSummaryDTO, 0, len(us))
    for _, u := range us {
        userList = append(userList, u.ToSummaryDTO())
    }

    out := usersOutput {
        Body: UserListSummaryDTO {
            Data:       userList,
            Page:       in.Page,
            PageSize:   in.PageSize,
            Total:      len(us),
        },
    }
	return &out, nil
}

// patch

func (h *Handler) handlePatchUser(ctx context.Context, in *PatchUserInput) (*userOutput, error) {
	updates := map[string]any{}
 	if err := populateUpdates(&updates, *in); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	u, err := h.updateUserFieldsByID(ctx, in.ID, updates)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &userOutput{Body: u.ToSummaryDTO()}, nil
}

func populateUpdates(updates *map[string]any, in PatchUserInput) error {
	if in.Body.Name != nil {
		if err := auth.ValidUserName(*in.Body.Name); err != nil {
			return err
		}
		(*updates)["name"] = *in.Body.Name
	}
	if in.Body.Email != nil {
		if err := auth.ValidUserEmail(*in.Body.Email); err != nil {
			return err
		}
		(*updates)["email"] = *in.Body.Email
	}
	return nil
}

// patch password

func (h *Handler) handlePatchPassword(ctx context.Context, in *PatchPasswordInput) (*userOutput, error) {
	return h.patchPassword(ctx, in)
}

func (h *Handler) patchPassword(ctx context.Context, in *PatchPasswordInput) (*userOutput, error) {
	u, err := h.getUserByID(ctx, in.ID)
    if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error404NotFound(err.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}

	match, err := auth.MatchPassword(in.Body.CurrentPassword, u.PasswordHash)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	if !match {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	if in.Body.NewPassword != in.Body.ConfirmPassword {
		return nil, huma.Error400BadRequest("new passwords do not match")
	}

	if err := auth.ValidUserPassword(in.Body.NewPassword); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	hash, err := auth.CreateHash(in.Body.NewPassword)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	u, err = h.updateUserFieldsByID(ctx, in.ID, map[string]any{"password_hash": hash})
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &userOutput{Body: u.ToSummaryDTO()}, nil
}

// delete

func (h *Handler) handleDeleteUser(ctx context.Context, in *deleteUserInput) (*userOutput, error) {
	err := h.deleteUserByID(ctx, in.ID)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
    return nil, nil
}

// get me

func (h *Handler) handleGetMe(ctx context.Context, in *struct{}) (*userOutput, error) {
	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	u, err := h.getUserByID(ctx, id)
    if errors.Is(err, errs.ErrNotFound) {
        return nil, huma.Error404NotFound(err.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
    return &userOutput{Body: u.ToSummaryDTO()}, nil
}

// patch me

func (h *Handler) handlePatchMe(ctx context.Context, in *PatchUserInput) (*userOutput, error) {
	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	updates := map[string]any{}
 	if err := populateUpdates(&updates, *in); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	u, err := h.updateUserFieldsByID(ctx, id, updates)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &userOutput{Body: u.ToSummaryDTO()}, nil
}

// patch password me

func (h *Handler) handlePatchPasswordMe(ctx context.Context, in *PatchPasswordInput) (*userOutput, error) {

	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	in.ID = id;
	return h.patchPassword(ctx, in)
}

// delete me

func (h *Handler) handleDeleteMe(ctx context.Context, in *deleteUserInput) (*userOutput, error) {
	id, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error404NotFound(errs.ErrNotFound.Error())
	}

	err = h.deleteUserByID(ctx, id)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, huma.Error404NotFound(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
    return nil, nil
}
