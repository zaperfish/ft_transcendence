package user

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"gorm.io/gorm"
)

func registerCreateUser(api huma.API, h dbHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/api/users",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusCreated,
	}, h.handleCreateUser)
}

func (h *dbHandler) handleCreateUser(ctx context.Context, in *createInput) (*userOutput, error) {
	u := User{
		Name: in.Body.Name,
	}

	err := gorm.G[User](h.db).Create(ctx, &u)
	if err != nil {
		return nil, err
	}

	return &userOutput{Body: u.toResponseDTO()}, nil
}

type createInput struct {
	Body createDTO
}

type createDTO struct {
	Name string `json:"name" maxLength:"30" example:"Max" doc:"username"`
}
