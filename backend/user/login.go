package user

import (
    // Std
	"context"
	"errors"
	"net/http"

    // Internal
	"ft_transcendence/backend/auth"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
    "gorm.io/gorm"
)

// login user
func registerLoginUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "login-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/login",
		Summary:		"Login",
		Description:	"Acquire a JWT cookie",
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

type LoginUserOutput struct {
	SetCookie http.Cookie 		`header:"Set-Cookie"`
    Body UserSummaryDTO
}


func (h *handler) handleLoginUser(ctx context.Context, in *loginUserInput) (*LoginUserOutput, error) {
    u, err := gorm.G[User](h.db).Where("name = ?", in.Body.Name).First(ctx)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, huma.Error401Unauthorized(gorm.ErrRecordNotFound.Error())
    }
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
	}
	
	match, err := auth.MatchPassword(in.Body.Password, u.PasswordHash)
	if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
    }
	if !match {
        return nil, huma.Error401Unauthorized(gorm.ErrRecordNotFound.Error())
	}

	cookie, err := auth.MakeJWTCookieFromID(u.ID)
    if err != nil {
        return nil, huma.Error500InternalServerError(err.Error())
    }

    out := &LoginUserOutput {
		SetCookie: cookie,
        Body: 	   u.ToSummaryDTO(),
    }

    return out, nil
}
