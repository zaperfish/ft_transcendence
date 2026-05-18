package user

import (
	"net/http"
	"context"

    "gorm.io/gorm"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func registerGet(api huma.API, h dbHandler) {
    huma.Register(api, huma.Operation{
        OperationID:    "get-user",
        Method:         http.MethodGet,
        Path:           "/api/users/{name}",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleGet)
}

func (h *dbHandler) handleGet(ctx context.Context, in *getInput) (*userOutput, error) {
    u, err := gorm.G[user](h.db).Where("name = ?", in.Name).First(ctx)
    if err != nil {
        return nil, err
    }
    return &userOutput{Body: u.toResponseDTO()}, nil
}

type getInput struct {
    Name string `path:"name" maxLength:"30" example:"Max" doc:"get user by name"`
}
