package chat

import (
	// Std
	"context"
	"errors"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	// Extern
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	clientSendBuffer               = 16
	clientWriteTimeout             = 10 * time.Second
	maxMessageCharacters           = 2000
	maxMessageBytes                = 8000
	maxWebSocketMessagePayloadSize = 12 * 1024
)

type Client struct {
	userID   uint
	userName string
	conn     *websocket.Conn
	// channel used as queue
	send chan MessageDTO
}

type messageCreator interface {
	createMessage(context.Context, *Message) error
}

func (c *Client) readLoop(ctx context.Context, room *Room, eventID uint, messages messageCreator) {
	for {
		var input createMessageInput
		if err := wsjson.Read(ctx, c.conn, &input); err != nil {
			return
		}

		message, err := newMessageFromInput(eventID, c.userID, input)
		if err != nil {
			_ = c.conn.Close(websocket.StatusPolicyViolation, err.Error())
			return
		}

		if err := messages.createMessage(ctx, &message); err != nil {
			log.Printf("failed to create chat message: event_id=%d user_id=%d err=%v", eventID, c.userID, err)
			_ = c.conn.Close(websocket.StatusInternalError, "could not create message")
			return
		}

		if !room.Broadcast(message.toDTO(c.userName)) {
			return
		}
	}
}

func (c *Client) writeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-c.send:
			if !ok {
				return
			}

			writeCtx, cancel := context.WithTimeout(ctx, clientWriteTimeout)
			err := wsjson.Write(writeCtx, c.conn, message)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

func newMessageFromInput(eventID uint, userID uint, input createMessageInput) (Message, error) {
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return Message{}, errors.New("message content cannot be empty")
	}
	if len(content) > maxMessageBytes {
		return Message{}, errors.New("message content exceeds 8000 bytes")
	}
	if utf8.RuneCountInString(content) > maxMessageCharacters {
		return Message{}, errors.New("message content exceeds 2000 characters")
	}

	return Message{
		EventID: eventID,
		UserID:  userID,
		Content: content,
	}, nil
}
