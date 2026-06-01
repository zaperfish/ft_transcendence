package chat

import "github.com/coder/websocket"

type Client struct {
	userID uint
	conn   *websocket.Conn
	send   chan Message
}
