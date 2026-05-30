package chat

type Room struct {
    eventID   uint
    // clients is used as a set of active WebSocket clients in this room
    // the bool value is not meaningful
    clients   map[*Client]bool
    join      chan *Client
    leave     chan *Client
    broadcast chan Message
}

func NewRoom(eventID uint) *Room {
    return &Room{
        eventID:   eventID,
        clients:   make(map[*Client]bool),
        join:      make(chan *Client),
        leave:     make(chan *Client),
        broadcast: make(chan Message),
    }
}
