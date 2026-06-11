package me

import (
    // Std
	"net/http"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func RegisterRoutes(api huma.API, h MeHandler) {

    huma.Register(api, huma.Operation{
        OperationID:    "get-me",
        Method:         http.MethodGet,
        Path:           "/api/me",
		Summary:		"Get logged in user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
		Security: 		[]map[string][]string{
							{"SessionToken": {}},
						},
    }, h.handleGetMe)

    huma.Register(api, huma.Operation{
        OperationID:    "patch-me",
        Method:         http.MethodPatch,
        Path:           "/api/me",
        Summary:        "Update logged in user",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
		Security: 		[]map[string][]string{
							{"SessionToken": {}},
						},
    }, h.handlePatchMe)

    huma.Register(api, huma.Operation{
        OperationID:    "patch-me-password",
        Method:         http.MethodPatch,
        Path:           "/api/me/password",
        Summary:        "Update logged in user's password",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
		Security: 		[]map[string][]string{
							{"SessionToken": {}},
						},
    }, h.handlePatchPasswordMe)

    huma.Register(api, huma.Operation{
        OperationID:    "delete-me",
        Method:         http.MethodDelete,
        Path:           "/api/me",
		Summary:		"Delete logged in user",
        DefaultStatus:  http.StatusNoContent,
        Tags:           []string{"Me"},
		Security: 		[]map[string][]string{
							{"SessionToken": {}},
						},
    }, h.handleDeleteMe)

		//   huma.Register(api, huma.Operation{
		//       OperationID:    "create-event-me",
		//       Method:         http.MethodDelete,
		//       Path:           "/api/me/create-event",
		// Summary:		"Create event as logged in user",
		//       DefaultStatus:  http.StatusNoContent,
		//       Tags:           []string{"Me"},
		//   }, h.handleCreateEventMe)

    huma.Register(api, huma.Operation{
        OperationID:    "join-event-me",
        Method:         http.MethodPost,
        Path:           "/api/me/join/{id}",
		Summary:		"Add logged in user to event",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
		Security: 		[]map[string][]string{
							{"SessionToken": {}},
						},
    }, h.handleJoinEventMe)

    huma.Register(api, huma.Operation{
        OperationID:    "leave-event-me",
        Method:         http.MethodDelete,
        Path:           "/api/me/leave/{id}",
		Summary:		"Remove logged in user from event",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Me"},
		Security: 		[]map[string][]string{
							{"SessionToken": {}},
						},
    }, h.handleLeaveEventMe)

	// huma.Register(api, huma.Operation{
	// 	OperationID:    "list-events-me",
	// 	Method:         http.MethodGet,
	// 	Path:           "/api/me/events",
	// 	Summary:		"List events logged in user is registered for",
	// 	DefaultStatus:  http.StatusOK,
	// 	Tags:           []string{"Me"},
	// 	Security: 		[]map[string][]string{
	// 						{"SessionToken": {}},
	// 					},
	// }, h.handleEventsMe)

		//   huma.Register(api, huma.Operation{
		//       OperationID:    "admin-events-me",
		//       Method:         http.MethodGet,
		//       Path:           "/api/me/admin-events",
		// Summary:		"Get events logged in user administers",
		//       DefaultStatus:  http.StatusOK,
		//       Tags:           []string{"Me"},
		//   }, h.handleAdminEventsMe)
}
