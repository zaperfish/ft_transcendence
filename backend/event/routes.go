package event

import (
	// Std
	"net/http"

	// Intern
	"ft_transcendence/backend/apikey"
	"ft_transcendence/backend/auth"

	// Extern

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(api huma.API, db *gorm.DB) {
	// Setup layers
	eventRepo := NewEventRepository(db)
	eventService := NewEventService(eventRepo, db)
	eventHandler := NewEventHandler(eventService)

	// Register POST /events
	huma.Register(api, huma.Operation{
		OperationID:   "create-event",
		Method:        http.MethodPost,
		Path:          "/api/events",
		Summary:       "Create event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusCreated,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
		// Security: []map[string][]string{
		// 	{"AdminPassword": {}},
		// },
	}, eventHandler.CreateEvent)

	// Register PATCH /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "update-event",
		Method:        http.MethodPatch,
		Path:          "/api/events/{id}",
		Summary:       "Update event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.UpdateEvent)

	// Register DELETE /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "delete-event",
		Method:        http.MethodDelete,
		Path:          "/api/events/{id}",
		Summary:       "Delete event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.DeleteEvent)

	// Register GET /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "get-event",
		Method:        http.MethodGet,
		Path:          "/api/events/{id}",
		Summary:       "Get event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.GetEvent)

	// Register GET /events
	huma.Register(api, huma.Operation{
		OperationID:   "list-events",
		Method:        http.MethodGet,
		Path:          "/api/events",
		Summary:       "List events",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.ListEvents)

	huma.Register(api, huma.Operation{
		OperationID:   "add-participant",
		Method:        http.MethodPost,
		Path:          "/api/events/{id}/participants",
		Summary:       "Add participant",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, eventHandler.AddParticipant)

	huma.Register(api, huma.Operation{
		OperationID:   "remove-participant",
		Method:        http.MethodDelete,
		Path:          "/api/events/{eventID}/participants/{userID}",
		Summary:       "Remove participant",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.RemoveParticipant)

	huma.Register(api, huma.Operation{
		OperationID:   "list-participants",
		Method:        http.MethodGet,
		Path:          "/api/events/{id}/participants",
		Summary:       "List participants",
		Tags:          []string{"Events", "Images"},
		DefaultStatus: http.StatusOK,
	}, eventHandler.ListParticipants)

	// images
	huma.Register(api, huma.Operation{
		OperationID:   "create-event-image",
		Method:        http.MethodPost,
		Path:          "/api/events/{id}/image",
		Summary:       "Create event image",
		Description:   "accepts image/jpeg and image/png",
		Tags:          []string{"Events", "Images"},
		DefaultStatus: http.StatusCreated,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.CreateImage)

	huma.Register(api, huma.Operation{
		OperationID:   "get-event-image",
		Method:        http.MethodGet,
		Path:          "/api/events/{id}/image",
		Summary:       "Get event image",
		Tags:          []string{"Events", "Images"},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.GetImage)

	huma.Register(api, huma.Operation{
		OperationID:   "update-event-image",
		Method:        http.MethodPatch,
		Path:          "/api/events/{id}/image",
		Summary:       "Update event image",
		Description:   "accepts image/jpeg and image/png",
		Tags:          []string{"Events", "Images"},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.UpdateImage)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-event-image",
		Method:        http.MethodDelete,
		Path:          "/api/events/{id}/image",
		Summary:       "Delete event image",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
		Middlewares:   huma.Middlewares{auth.Verifier(api), auth.Refresher(api)},
	}, eventHandler.DeleteImage)

	// public api
	v1 := huma.NewGroup(api, "/api/v1")
	v1.UseMiddleware(apikey.ApiKeyVerifier(api, db))

	// Register POST /events
	huma.Register(v1, huma.Operation{
		OperationID:   "v1-create-event",
		Method:        http.MethodPost,
		Path:          "/events",
		Summary:       "Create event",
		Tags:          []string{"Public Events"},
		DefaultStatus: http.StatusCreated,
		Security: []map[string][]string{
			{"ApiKey": {}},
		},
	}, eventHandler.CreateEvent)

	// Register PATCH /events/{id}
	huma.Register(v1, huma.Operation{
		OperationID:   "v1-update-event",
		Method:        http.MethodPatch,
		Path:          "/events/{id}",
		Summary:       "Update event",
		Tags:          []string{"Public Events"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"ApiKey": {}},
		},
	}, eventHandler.UpdateEvent)

	// Register DELETE /events/{id}
	huma.Register(v1, huma.Operation{
		OperationID:   "v1-delete-event",
		Method:        http.MethodDelete,
		Path:          "/events/{id}",
		Summary:       "Delete event",
		Tags:          []string{"Public Events"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"ApiKey": {}},
		},
	}, eventHandler.DeleteEvent)

	// Register GET /events/{id}
	huma.Register(v1, huma.Operation{
		OperationID:   "v1-get-event",
		Method:        http.MethodGet,
		Path:          "/events/{id}",
		Summary:       "Get event",
		Tags:          []string{"Public Events"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"ApiKey": {}},
		},
	}, eventHandler.GetEvent)

	// Register GET /events
	huma.Register(v1, huma.Operation{
		OperationID:   "v1-list-events",
		Method:        http.MethodGet,
		Path:          "/events",
		Summary:       "List events",
		Tags:          []string{"Public Events"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"ApiKey": {}},
		},
	}, eventHandler.ListEvents)

	huma.Register(v1, huma.Operation{
		OperationID:   "v1-add-participant",
		Method:        http.MethodPost,
		Path:          "/events/{id}/participants",
		Summary:       "Add participant",
		Tags:          []string{"Public Events"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"ApiKey": {}},
		},
	}, eventHandler.AddParticipant)

	huma.Register(v1, huma.Operation{
		OperationID:   "v1-remove-participant",
		Method:        http.MethodDelete,
		Path:          "/events/{eventID}/participants/{userID}",
		Summary:       "Remove participant",
		Tags:          []string{"Public Events"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"ApiKey": {}},
		},
	}, eventHandler.RemoveParticipant)

	huma.Register(v1, huma.Operation{
		OperationID:   "v1-list-participants",
		Method:        http.MethodGet,
		Path:          "/events/{id}/participants",
		Summary:       "List participants",
		Tags:          []string{"Public Events"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"ApiKey": {}},
		},
	}, eventHandler.ListParticipants)
}
