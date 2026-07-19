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
	Content string `gorm:"not null;check:chk_messages_content_length,char_length(content) >= 1 AND char_length(content) <= 2000 AND octet_length(content) <= 8000"`
}

type MessageDTO struct {
	ID         uint      `json:"id" doc:"ID of the message"`
	EventID    uint      `json:"event_id" doc:"ID of the event"`
	UserID     uint      `json:"user_id" doc:"ID of the sender"`
	SenderName string    `json:"sender_name" doc:"Username of the sender"`
	Content    string    `json:"content" doc:"Message content"`
	CreatedAt  time.Time `json:"created_at" doc:"Time the message was created"`
}

func (m *Message) toDTO(senderName string) MessageDTO {
	return MessageDTO{
		ID:         m.ID,
		EventID:    m.EventID,
		UserID:     m.UserID,
		SenderName: senderName,
		Content:    m.Content,
		CreatedAt:  m.CreatedAt,
	}
}

func messagesToDTOsOldestFirst(messages []Message, senderNames map[uint]string) []MessageDTO {
	dtos := make([]MessageDTO, len(messages))
	for i, message := range messages {
		dtos[len(messages)-1-i] = message.toDTO(senderNames[message.UserID])
	}

	return dtos
}
