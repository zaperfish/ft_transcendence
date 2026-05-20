package user

import (
    // Std
	"context"
	"net/http"

    // Internal
	// "ft_transcendence/backend/app"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
    // "github.com/go-chi/jwtauth/v5"
    "gorm.io/gorm"
)

type userLoginDTO struct {
    Name string     `json:"name"`
    Password string `json:"password"`
}

func registerLoginUser(api huma.API, h Handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "login-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/login",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleLoginUser)
}

type loginUserInput struct {
    Body userLoginDTO
}

type UserLoginResponseDTO struct {
    ResponseDTO UserResponseDTO
    AccessToken string              `json:"access_token"`
}

type loginUserOutput struct {
    Body UserLoginResponseDTO
}

// type UserLoginResponseDTO struct {
//     AccessToken string              `json:"access_token"`
//     TokenType string                `json:"token_type"`
// }
//
// type loginUserOutput struct {
//     Body UserLoginResponseDTO
// }

func (h *Handler) handleLoginUser(ctx context.Context, in *loginUserInput) (*loginUserOutput, error) {
    u, err := gorm.G[user](h.app.DB).Where("name = ?", in.Body.Name).First(ctx)
    if err != nil {
        return nil, err
    }

    if u.Password != in.Body.Password {
        return nil, gorm.ErrRecordNotFound
    }

    _, t, err := h.app.TokenAuth.Encode(map[string]any{"user_id": u.ID})
    if err != nil {
        return nil, err
    }
    
    out := &loginUserOutput {
        Body: UserLoginResponseDTO {
            ResponseDTO:    u.toResponseDTO(),
            AccessToken:    t,
        },
    }
    return out, nil
}
