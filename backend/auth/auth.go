package auth

import (
    // Std
	"context"
	"net/http"

    // Internal
	"ft_transcendence/backend/app"
	"ft_transcendence/backend/user"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
    // "github.com/go-chi/jwtauth/v5"
    "gorm.io/gorm"
)

func RegisterApi(api huma.API, app *app.App) {
    app.DB.AutoMigrate(&user.User{})

    h := Handler {app: app}
    registerRegisterUser(api, h);
    registerLoginUser(api, h);
}

type Handler struct {
    app *app.App
}

// login user
func registerRegisterUser(api huma.API, h Handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "register-user",
        Method:         http.MethodPost,
        Path:           "/api/register",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Authentification"},
    }, h.handleCreateUser)
}

func (h *Handler) handleCreateUser(ctx context.Context, in *createInput) (*userOutput, error) {
    u := user.User {
        Name:       in.Body.Name,
        Password:   in.Body.Password,
    }

    err := gorm.G[user.User](h.app.DB).Create(ctx, &u)
    if err != nil {
        return nil, err
    }

    return &userOutput{Body: u.ToResponseDTO()}, nil
}

type createInput struct {
    Body CreateDTO
}

type CreateDTO struct {
    Name string     `json:"name" maxLength:"30" example:"Max" doc:"username"`
    Password string `json:"password" example:"secret" doc:"password"`
}

type userOutput struct {
    Body user.UserResponseDTO
}

// login user
func registerLoginUser(api huma.API, h Handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "login-user",
        Method:         http.MethodPost,
        Path:           "/api/login",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Authentification"},
    }, h.handleLoginUser)
}

type userLoginDTO struct {
    Name string     `json:"name" example:"Max"`
    Password string `json:"password" example:"secret"`
}

type loginUserInput struct {
    Body userLoginDTO
}

type UserLoginResponseDTO struct {
    ResponseDTO user.UserResponseDTO
    AccessToken string              `json:"access_token"`
}

type loginUserOutput struct {
    Body UserLoginResponseDTO
}

func (h *Handler) handleLoginUser(ctx context.Context, in *loginUserInput) (*loginUserOutput, error) {
    u, err := gorm.G[user.User](h.app.DB).Where("name = ?", in.Body.Name).First(ctx)
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
            ResponseDTO:    u.ToResponseDTO(),
            AccessToken:    t,
        },
    }
    return out, nil
}
