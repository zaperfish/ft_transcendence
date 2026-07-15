package user

import (
    // Std
	"net/http"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func RegisterPublicRoutes(api huma.API, h UserHandler) {

    huma.Register(api, huma.Operation{
        OperationID:    "register-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/register",
		Summary:		"Register a new user",
        DefaultStatus:  http.StatusCreated,
        Tags:           []string{"Authentification"},
    }, h.handleRegisterUser)

    huma.Register(api, huma.Operation{
        OperationID:    "login-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/login",
		Summary:		"Login",
		Description:	"Acquire a testing JWT",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Authentification"},
    }, h.handleLoginUser)

    huma.Register(api, huma.Operation{
        OperationID:    "logout-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/logout",
		Summary:		"Logout",
		Description:	"Instructs browser to delete JWT cookie",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Authentification"},
    }, h.handleLogoutUser)

    huma.Register(api, huma.Operation{
        OperationID:    "get-session-token",
        Method:         http.MethodGet,
        Path:           "/api/auth/token",
		Summary:		"Get a JWT session token",
		Description:	"Creates user \"dummy\" \"dummy@dummy.com\" and/or returns a jwt session token for them (this endpoint should be removed in production)",
        DefaultStatus:  http.StatusCreated,
        Tags:           []string{"Authentification"},
    }, h.handleGetToken)
}

func RegisterProtectedRoutes(api huma.API, h UserHandler) {



    huma.Register(api, huma.Operation{
        OperationID:    "get-user-by-id",
        Method:         http.MethodGet,
        Path:           "/api/users/{id}",
		Summary:		"Get a user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleGetUser)

    huma.Register(api, huma.Operation{
        OperationID:    "get-users",
        Method:         http.MethodGet,
        Path:           "/api/users",
		Summary:		"Get a list of user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleGetUsers)
}
