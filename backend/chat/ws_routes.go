package chat

import (
	// Extern
	"github.com/go-chi/chi/v5"
)

func RegisterWebSocketRoutes(r chi.Router, h Handler) {
	r.Get("/api/events/{id}/chat/ws", h.handleEventChatWebSocket)
}
