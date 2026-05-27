package user

import (
    // Std
	"context"
    "errors"
	"net/http"

	// Internal
	"ft_transcendence/backend/auth"

    // External
    "gorm.io/gorm"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func registerPatchUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "patch-user",
        Method:         http.MethodPatch,
        Path:           "/api/users/{id}",
        Summary:        "Update user by ID",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handlePatchUser)
}

func registerPatchPassword(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "patch-user-password",
        Method:         http.MethodPatch,
        Path:           "/api/users/{id}/password",
        Summary:        "Update a user's password by ID",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handlePatchPassword)
}

type PatchUserDTO struct {
    Name *string     `json:"name,omitempty" maxLength:"30" example:"Max" doc:"username"`
    Email *string    `json:"email,omitempty" example:"max@email.com" doc:"email address"`
}

type PatchPasswordDTO struct {
    NewPassword string `json:"newpassword" example:"newsecret" doc:"new password"`
    ConfirmPassword string `json:"confirm_password" example:"newsecret" doc:"confirm password"`
    CurrentPassword string `json:"current_password" example:"secret" doc:"current password"`
}

type PatchPasswordInput struct {
	ID uint	`path:"id" doc:"User ID" example:"1"`
	Body PatchPasswordDTO
}

type PatchUserInput struct {
	ID uint	`path:"id" doc:"User ID" example:"1"`
	Body PatchUserDTO
}

func (h *handler) handlePatchPassword(ctx context.Context, in *PatchPasswordInput) (*userOutput, error) {
	// sub, err := auth.GetSubClaim(ctx)
	// if err != nil || sub != strconv.FormatUint(uint64(in.ID), 10) {
	// 	return nil, huma.Error401Unauthorized("wrong permissions")
	// }

	u, err := h.getByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	match, err := auth.MatchPassword(in.Body.CurrentPassword, u.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("old password does not match")
	}

	if in.Body.NewPassword != in.Body.ConfirmPassword {
		return nil, errors.New("new passwords do not match")
	}

	if err := auth.ValidUserPassword(in.Body.NewPassword); err != nil {
		return nil, err
	}

	hash, err := auth.CreateHash(in.Body.NewPassword)
	if err != nil {
		return nil, huma.Error500InternalServerError("")
	}

	u, err = h.updateFieldsByID(ctx, in.ID, map[string]any{"password_hash": hash})
	if err != nil {
		return nil, err
	}

	return &userOutput{Body: u.ToSummaryDTO()}, nil
}

func (h *handler) handlePatchUser(ctx context.Context, in *PatchUserInput) (*userOutput, error) {
	updates := map[string]any{}
 	if err := populateUpdates(&updates, *in, h.db, ctx); err != nil {
		return nil, err
	}

	u, err := h.updateFieldsByID(ctx, in.ID, updates)
	if err != nil {
		return nil, err
	}

	return &userOutput{Body: u.ToSummaryDTO()}, nil
}

func populateUpdates(updates *map[string]any, in PatchUserInput, db *gorm.DB, ctx context.Context) error {
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
