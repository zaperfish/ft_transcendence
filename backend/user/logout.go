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
)

// logout user

func registerLogoutUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "logout-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/logout",
		Summary:		"Logout",
		Description:	"Instructs browser to delete JWT cookie",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Authentification"},
    }, h.handleLogoutUser)
}

type LogoutUserOutput struct {
	SetCookie http.Cookie 		`header:"Set-Cookie"`
}

func (h *handler) handleLogoutUser(ctx context.Context, in *struct{}) (*LogoutUserOutput, error) {

    out := &LogoutUserOutput {
		SetCookie: auth.MakeLogoutCookie(),
    }

    return out, nil
}
