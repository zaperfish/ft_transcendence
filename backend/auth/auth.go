package auth

import (
    // Std
	"context"
	"fmt"
	"net/http"
	"time"

    // Internal
	"ft_transcendence/backend/app"
	"ft_transcendence/backend/user"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
    // "github.com/go-chi/jwtauth/v5"
    "github.com/go-chi/jwtauth/v5"
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
    Email string    `json:"email" example:"max@email.com" doc:"email address"`
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

type LoginUserOutput struct {
	SetCookie http.Cookie 		`header:"Set-Cookie"`
    Body UserLoginResponseDTO
}

// func refreshJWTCookie(ctx huma.Context, next func(ctx huma.Context)) {
// 	_, claims, err := jwtauth.FromContext(ctx.Context())
// 	if err != nil {
// 		return
// 	}
//
// 	id, ok := claims["user_id"].(uint)
// 	if !ok {
// 		return
// 	}
//
// 	iat, ok := claims["iat"].(time.Time)
// 	if !ok {
// 		return
// 	}
//
// 	claims = map[string]any {
// 		"user_id":		id,
// 		"exp":			time.Now().Add(30 * time.Minute).Unix(),
// 		"iat":			iat,
// 	}
// }

func makeJWT(tokenAuth *jwtauth.JWTAuth, uid uint) (string, error) {
	claims := map[string]any {
		"user_id":		uid,
		"exp":			time.Now().Add(30 * time.Minute).Unix(),
		"iat":			time.Now().Unix(),
	}
    _, t, err := tokenAuth.Encode(claims)
    if err != nil {
        return "", err
    }
	return t, nil
}

func makeJWTCookie(tokenAuth *jwtauth.JWTAuth, uid uint) (http.Cookie, error) {
	t, err := makeJWT(tokenAuth, uid)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie {
		Name:		"auth_token",
		Value:		t,
		Path:		"/",
		Expires:	time.Now().Add(15 * time.Minute),
		HttpOnly:	true,
		Secure:		true,
		SameSite:	http.SameSiteNoneMode,
	}, nil
}

func (h *Handler) handleLoginUser(ctx context.Context, in *loginUserInput) (*LoginUserOutput, error) {
    u, err := gorm.G[user.User](h.app.DB).Where("name = ?", in.Body.Name).First(ctx)
    if err != nil {
        return nil, err
    }

    if u.Password != in.Body.Password {
        return nil, gorm.ErrRecordNotFound
    }

	cookie, err := makeJWTCookie(h.app.TokenAuth, u.ID)
    if err != nil {
        return nil, err
    }

    out := &LoginUserOutput {
		SetCookie: cookie,
        Body: UserLoginResponseDTO {
            ResponseDTO:    u.ToResponseDTO(),
        },
    }
	fmt.Println(out)

    return out, nil
}
