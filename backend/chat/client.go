package chat

import "github.com/coder/websocket"

const clientSendBuffer = 16

type Client struct {
	userID uint
	conn   *websocket.Conn
	send   chan Message
}
