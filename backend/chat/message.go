package chat

import (
	// Std
	"time"

	// Extern
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model

	EventID uint   `gorm:"not null;index"`
	UserID  uint   `gorm:"not null;index"`
	Content string `gorm:"not null;check:length(content) >= 1"`
}

type MessageDTO struct {
	ID        uint      `json:"id" doc:"ID of the message"`
	EventID   uint      `json:"event_id" doc:"ID of the event"`
	UserID    uint      `json:"user_id" doc:"ID of the sender"`
	Content   string    `json:"content" doc:"Message content"`
	CreatedAt time.Time `json:"created_at" doc:"Time the message was created"`
}

func (m *Message) toDTO() MessageDTO {
	return MessageDTO{
		ID:        m.ID,
		EventID:   m.EventID,
		UserID:    m.UserID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
	}
}

func messagesToDTOsOldestFirst(messages []Message) []MessageDTO {
	dtos := make([]MessageDTO, len(messages))
	for i, message := range messages {
		dtos[len(messages)-1-i] = message.toDTO()
	}

	return dtos
}
