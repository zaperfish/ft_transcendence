package auth

import (
    // Std
	"context"
	"net/http"
	"strconv"
	"time"

    // Internal
	"ft_transcendence/backend/user"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
    "gorm.io/gorm"
)

func RegisterApi(api huma.API, db *gorm.DB ) {
    db.AutoMigrate(&user.User{})

	h := handler{db: db}
    registerRegisterUser(api, h);
    registerLoginUser(api, h);
    registerLogoutUser(api, h);
}

type handler struct {
    db *gorm.DB
}

// login user
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
    u := user.User {
        Name:       in.Body.Name,
        Email:      in.Body.Email,
        Password:   in.Body.Password,
    }

	err := user.ValidateUser(u)
	if err != nil {
		return nil, err
	}

    err = gorm.G[user.User](h.db).Create(ctx, &u)
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

type userOutput struct {
    Body user.UserSummaryDTO
}

// login user
func registerLoginUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "login-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/login",
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
    Body user.UserSummaryDTO
}

func makeJWT(sub string) (string, error) {
	claims := map[string]any {
		"sub":		sub,
		"exp":		time.Now().Add(jwtExpirationTime).Unix(),
		"iat":		time.Now().Unix(),
	}
    _, ts, err := tokenAuth.Encode(claims)
    if err != nil {
        return "", err
    }
	return ts, nil
}

func makeJWTCookie(sub string) (http.Cookie, error) {
	t, err := makeJWT(sub)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie {
		Name:		"auth_token",
		Value:		t,
		Path:		"/",
		Expires:	time.Now().Add(jwtExpirationTime),
		HttpOnly:	true,
		Secure:		true,
		SameSite:	http.SameSiteLaxMode,
	}, nil
}

func (h *handler) handleLoginUser(ctx context.Context, in *loginUserInput) (*LoginUserOutput, error) {
    u, err := gorm.G[user.User](h.db).Where("name = ?", in.Body.Name).First(ctx)
    if err != nil {
        return nil, err
    }

    if u.Password != in.Body.Password {
        return nil, gorm.ErrRecordNotFound
    }

	cookie, err := makeJWTCookie(strconv.FormatUint(uint64(u.ID), 10))
    if err != nil {
        return nil, err
    }

    out := &LoginUserOutput {
		SetCookie: cookie,
        Body: 	   u.ToSummaryDTO(),
    }

    return out, nil
}

// logout user

func registerLogoutUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "logout-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/logout",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Authentification"},
    }, h.handleLogoutUser)
}

func makeJWTDeleteCookie() (http.Cookie, error) {
	return http.Cookie {
	}, nil
}

type LogoutUserOutput struct {
	SetCookie http.Cookie 		`header:"Set-Cookie"`
}

func (h *handler) handleLogoutUser(ctx context.Context, in *struct{}) (*LogoutUserOutput, error) {

    out := &LogoutUserOutput {
		SetCookie: http.Cookie {
			Name:		"auth_token",
			Value:		"",
			Path:		"/",
			HttpOnly:	true,
			Secure:		true,
			MaxAge:		-1,
		},
    }

    return out, nil
}
