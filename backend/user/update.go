package user

import (
	"context"
    "fmt"
	"net/http"

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
    Password *string `json:"password,omitempty" example:"secret" doc:"password"`
}

type PatchUserInput struct {
	ID uint	`path:"id" doc:"User ID"`
	Body PatchUserDTO
}

func (h *handler) handlePatchUser(ctx context.Context, in *PatchUserInput) (*userOutput, error) {
	claims := ctx.Value("claims").(map[string]any)

	// this still looks dodgey, change later
	if uint(claims["user_id"].(float64)) != in.ID {
		return nil, huma.Error401Unauthorized("wrong permissions")
	}

	updates := map[string]any{}

	if in.Body.Name != nil {
		updates["name"] = *in.Body.Name
	}
	if in.Body.Email != nil {
		updates["email"] = *in.Body.Email
	}
	if in.Body.Password != nil {
		updates["password"] = *in.Body.Password
	}

	_, err := gorm.G[map[string]any](h.db.Debug()).Table("users").Where("id = ?", in.ID).Updates(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to save patched user: %w", err)
	}

	updated, err := gorm.G[User](h.db.Debug()).Where("id = ?", in.ID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated user: %w", err)
	}

	return &userOutput{Body: updated.ToResponseDTO()}, nil
}
