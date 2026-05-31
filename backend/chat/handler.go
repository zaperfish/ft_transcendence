package chat

import (
	// Std
	"context"

	// Intern
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/event"

	// Extern
	"github.com/danielgtaylor/huma/v2"
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
