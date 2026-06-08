package chat

import (
	// Std
	"context"
	"errors"
	"log"
	"strings"
	"time"

	// Extern
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	clientSendBuffer   = 16
	clientWriteTimeout = 10 * time.Second
)

type Client struct {
	userID uint
	conn   *websocket.Conn
	// channel used as queue
	send   chan Message
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

		if !room.Broadcast(message) {
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
			err := wsjson.Write(writeCtx, c.conn, message.toDTO())
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

	return Message{
		EventID: eventID,
		UserID:  userID,
		Content: content,
	}, nil
}
