package chat

import (
	// Std
	"context"
	"net/http"
	"strconv"

	// Intern
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/event"

	// Extern
	"github.com/coder/websocket"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

const messageHistoryLimit = 50

type Handler struct {
	DB  *gorm.DB
	Hub *Hub
}

func NewHandler(db *gorm.DB) Handler {
	return Handler{
		DB:  db,
		Hub: NewHub(),
	}
}

func (h *Handler) handleGetEventMessages(ctx context.Context, input *getMessagesInput) (*messagesOutput, error) {
	userID, err := auth.UidFromCtx(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized(err.Error())
	}

	eventExists, err := event.EventExists(ctx, h.DB, input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	if !eventExists {
		return nil, huma.Error404NotFound("event not found")
	}

	isParticipant, err := event.IsParticipant(ctx, h.DB, input.ID, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	if !isParticipant {
		return nil, huma.Error403Forbidden("user is not a participant")
	}

	messages, err := h.getMessagesByEventID(ctx, input.ID, MessageFilter{
		BeforeID: input.BeforeID,
		Limit:    messageHistoryLimit,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &messagesOutput{
		Body: MessageListDTO{
			Data: messagesToDTOsOldestFirst(messages),
		},
	}, nil
}

// raw HTTP/Chi handler for websockets
func (h *Handler) handleEventChatWebSocket(w http.ResponseWriter, r *http.Request) {
	eventIDParam := chi.URLParam(r, "id")

	eventID64, err := strconv.ParseUint(eventIDParam, 10, strconv.IntSize)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	eventID := uint(eventID64)

	userID, err := auth.UidFromRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventExists, err := event.EventExists(r.Context(), h.DB, eventID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !eventExists {
		http.Error(w, "event not found", http.StatusNotFound)
		return
	}

	isParticipant, err := event.IsParticipant(r.Context(), h.DB, eventID, userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	client := &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan Message),
	}

	_ = eventID
	_ = client
}
