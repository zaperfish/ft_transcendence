package chat

import (
	// Std
	"context"

	// Intern
	"ft_transcendence/backend/errs"

	// Extern
	"gorm.io/gorm"
)

type MessageFilter struct {
	BeforeID uint
	Limit    int
}

func normalizeMessageLimit(limit int) int {
	if limit <= 0 || limit > messageHistoryLimit {
		return messageHistoryLimit
	}

	return limit
}

func (h Handler) createMessage(ctx context.Context, message *Message) error {
	err := gorm.G[Message](h.DB).Create(ctx, message)
	return errs.ErrorDB(err)
}

func (h Handler) getMessagesByEventID(ctx context.Context, eventID uint, filter MessageFilter) ([]Message, error) {
	query := gorm.G[Message](h.DB).
		Where("event_id = ?", eventID).
		Order("id DESC").
		Limit(normalizeMessageLimit(filter.Limit))

	if filter.BeforeID != 0 {
		query = query.Where("id < ?", filter.BeforeID)
	}

	messages, err := query.Find(ctx)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	return messages, nil
}
