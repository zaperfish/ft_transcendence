package user

import (
    // Std
	"context"
    "fmt"
	"net/http"
	"strconv"

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
        Summary:        "Update user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handlePatchUser)
}

type PatchUserDTO struct {
    Name *string     `json:"name,omitempty" maxLength:"30" example:"Max" doc:"username"`
    Email *string    `json:"email,omitempty" example:"max@email.com" doc:"email address"`
    Password *string `json:"password,omitempty" example:"newsecret" doc:"password"`
    OldPassword *string `json:"old_password,omitempty" example:"secret" doc:"old password"`
}

type PatchUserInput struct {
	ID uint	`path:"id" doc:"User ID"`
	Body PatchUserDTO
}

func (h *handler) handlePatchUser(ctx context.Context, in *PatchUserInput) (*userOutput, error) {
	sub, err := auth.GetSubClaim(ctx)
	if err != nil || sub != strconv.FormatUint(uint64(in.ID), 10) {
		return nil, huma.Error401Unauthorized("wrong permissions")
	}

	updates := map[string]any{}
 	if err := populateUpdates(&updates, *in, h.db, ctx); err != nil {
		return nil, err
	}

	_, err = gorm.G[map[string]any](h.db.Debug()).Table("users").Where("id = ?", in.ID).Updates(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to save patched user: %w", err)
	}

	updated, err := gorm.G[User](h.db.Debug()).Where("id = ?", in.ID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated user: %w", err)
	}

	return &userOutput{Body: updated.ToSummaryDTO()}, nil
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
 	if err := populatePasswordUpdate(updates, in, db, ctx); err != nil {
		return err
	}
	return nil
}

func populatePasswordUpdate(updates *map[string]any, in PatchUserInput, db *gorm.DB, ctx context.Context) error {
	if in.Body.Password != nil {
		if in.Body.OldPassword == nil {
			return gorm.ErrRecordNotFound
		}
		u, err := gorm.G[User](db).Where("id = ?", in.ID).First(ctx)
		if err != nil {
			return err
		}
		match, err := auth.MatchPassword(*in.Body.OldPassword, u.PasswordHash)
		if err != nil {
			return huma.Error500InternalServerError("")
		}
		if !match {
			return gorm.ErrRecordNotFound
		}
		if err := auth.ValidUserPassword(*in.Body.Password); err != nil {
			return err
		}
		hash, err := auth.CreateHash(*in.Body.Password)
		if err != nil {
			return huma.Error500InternalServerError("")
		}
		(*updates)["password_hash"] = hash
	}
	return nil
}
