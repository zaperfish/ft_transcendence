package user

import (
    // Std
	"net/http"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func RegisterPublicRoutes(api huma.API, h Handler) {

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
		Description:	"Acquire a JWT cookie",
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
}

func RegisterProtectedRoutes(api huma.API, h Handler) {

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

    huma.Register(api, huma.Operation{
        OperationID:    "patch-user",
        Method:         http.MethodPatch,
        Path:           "/api/users/{id}",
        Summary:        "Update a user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handlePatchUser)

    huma.Register(api, huma.Operation{
        OperationID:    "patch-user-password",
        Method:         http.MethodPatch,
        Path:           "/api/users/{id}/password",
        Summary:        "Update a user's password",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handlePatchPassword)

    huma.Register(api, huma.Operation{
        OperationID:    "delete-user-by-id",
        Method:         http.MethodDelete,
        Path:           "/api/users/{id}",
		Summary:		"Delete a user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleDeleteUser)

	// me
    huma.Register(api, huma.Operation{
        OperationID:    "get-me",
        Method:         http.MethodGet,
        Path:           "/api/me",
		Summary:		"Get logged in user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
    }, h.handleGetMe)

    huma.Register(api, huma.Operation{
        OperationID:    "patch-me",
        Method:         http.MethodPatch,
        Path:           "/api/me",
        Summary:        "Update logged in user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
    }, h.handlePatchMe)

    huma.Register(api, huma.Operation{
        OperationID:    "patch-me-password",
        Method:         http.MethodPatch,
        Path:           "/api/me/password",
        Summary:        "Update logged in user's password",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
    }, h.handlePatchPasswordMe)

    huma.Register(api, huma.Operation{
        OperationID:    "delete-me",
        Method:         http.MethodDelete,
        Path:           "/api/me",
		Summary:		"Delete logged in user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
    }, h.handleDeleteMe)
}
