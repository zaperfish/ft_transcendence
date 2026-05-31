package chat

type Room struct {
	eventID uint
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

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					// prioritizes server stability
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
