package user

import (
    // Std
	"context"
	"net/http"

    // Internal
	"ft_transcendence/backend/auth"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
    "gorm.io/gorm"
)

// register user
func registerRegisterUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "register-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/register",
        DefaultStatus:  http.StatusCreated,
        Tags:           []string{"Authentification"},
    }, h.handleCreateUser)
}

func (h *handler) handleCreateUser(ctx context.Context, in *createInput) (*userOutput, error) {

	if err := validateParameters(&in.Body); err != nil {
		return nil, err
	}

	hash, err := auth.CreateHash(in.Body.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError("")
	}

    u := User {
        Name:       	in.Body.Name,
        Email:      	in.Body.Email,
        PasswordHash:   hash,
    }

    err = gorm.G[User](h.db).Create(ctx, &u)
    if err != nil {
        return nil, err
    }

    return &userOutput{Body: u.ToSummaryDTO()}, nil
}

type createInput struct {
    Body CreateDTO
}

type CreateDTO struct {
    Name string     `json:"name" maxLength:"30" example:"Max" doc:"username"`
    Email string    `json:"email" example:"max@email.com" doc:"email address"`
    Password string `json:"password" example:"secret" doc:"password"`
}

func validateParameters(u *CreateDTO) error {
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
