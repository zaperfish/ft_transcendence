package chat

import (
	// Std
	"context"

	// Intern
	"ft_transcendence/backend/errs"
	"ft_transcendence/backend/user"

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

func (h Handler) getMessageSenderNames(ctx context.Context, messages []Message) (map[uint]string, error) {
	if len(messages) == 0 {
		return map[uint]string{}, nil
	}

	userIDs := make([]uint, 0, len(messages))
	seen := make(map[uint]struct{})
	for _, message := range messages {
		if _, ok := seen[message.UserID]; ok {
			continue
		}
		seen[message.UserID] = struct{}{}
		userIDs = append(userIDs, message.UserID)
	}

	return h.getUserNamesByIDs(ctx, userIDs)
}

func (h Handler) getUserNamesByIDs(ctx context.Context, userIDs []uint) (map[uint]string, error) {
	senderNames := make(map[uint]string)
	if len(userIDs) == 0 {
		return senderNames, nil
	}

	users, err := gorm.G[user.User](h.DB).
		Select("id", "name").
		Where("id IN ?", userIDs).
		Find(ctx)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	for _, sender := range users {
		senderNames[sender.ID] = sender.Name // build lookup map
	}

	return senderNames, nil
}
