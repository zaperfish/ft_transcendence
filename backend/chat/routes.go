package chat

import (
	// Std
	"net/http"

	// Extern
	"github.com/danielgtaylor/huma/v2"
)

func RegisterProtectedRoutes(api huma.API, h Handler) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-event-chat-messages",
		Method:        http.MethodGet,
		Path:          "/api/events/{id}/chat/messages",
		Summary:       "Get event chat messages",
		Tags:          []string{"Chat"},
		DefaultStatus: http.StatusOK,
	}, h.handleGetEventMessages)
}
