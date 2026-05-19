package user

import (
	"net/http"
	"context"

    "gorm.io/gorm"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func registerCreateUser(api huma.API, h dbHandler) {
    huma.Register(api, huma.Operation{
        OperationID:    "create-user",
        Method:         http.MethodPost,
        Path:           "/api/users",
        Tags:           []string{"Users"},
        DefaultStatus:  http.StatusCreated,
    }, h.handleCreateUser)
}

func (h *dbHandler) handleCreateUser(ctx context.Context, in *createInput) (*userOutput, error) {
    u := user {
        Name:       in.Body.Name,
        Password:   in.Body.Password,
    }

    err := gorm.G[user](h.db).Create(ctx, &u)
    if err != nil {
        return nil, err
    }

    return &userOutput{Body: u.toResponseDTO()}, nil
}

type createInput struct {
    Body createDTO
}

type createDTO struct {
    Name string     `json:"name" maxLength:"30" example:"Max" doc:"username"`
    Password string `json:"password" example:"secret" doc:"password"`
}
