package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"gorm.io/gorm"
)

// registerGetUser

func registerDeleteUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "delete-user-by-id",
        Method:         http.MethodDelete,
        Path:           "/api/users/{id}",
		Summary:		"Delete a user by ID",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleDeleteUser)
}

func (h *handler) handleDeleteUser(ctx context.Context, in *deleteUserInput) (*userOutput, error) {
    rows, err := gorm.G[User](h.db).Where("id = ?", in.ID).Delete(ctx)
    if err != nil {
        return nil, err
    }
	if rows == 0 {
		return nil, errors.New("no user deleted")
	}
    return nil, nil
}

type deleteUserInput struct {
	ID uint	`path:"id" doc:"User ID" example:"1"`
}
