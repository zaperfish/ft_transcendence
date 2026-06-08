package chat

type Room struct {
	eventID uint
	// clients is used as a set of active WebSocket clients in this room
	// the bool value is not meaningful
	clients   map[*Client]bool
	// channels used as event queues
	join      chan *Client
	leave     chan *Client
	broadcast chan Message
	// channel used as shutdown signal
	done      chan struct{}
	onEmpty   func(uint, *Room)
}

func NewRoom(eventID uint) *Room {
	return newRoom(eventID, nil)
}

func newRoom(eventID uint, onEmpty func(uint, *Room)) *Room {
	return &Room{
		eventID:   eventID,
		clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan Message),
		done:      make(chan struct{}),
		onEmpty:   onEmpty,
	}
}

func (r *Room) Join(client *Client) bool {
	select {
	case r.join <- client:
		return true
	case <-r.done:
		return false
	}
}

func (r *Room) Leave(client *Client) {
	select {
	case r.leave <- client:
	case <-r.done:
	}
}

func (r *Room) Broadcast(message Message) bool {
	select {
	case r.broadcast <- message:
		return true
	case <-r.done:
		return false
	}
}

func (r *Room) run() {
	defer close(r.done)

	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			if len(r.clients) == 0 {
				r.closeIfEmpty()
				return
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
			if len(r.clients) == 0 {
				r.closeIfEmpty()
				return
			}
		}
	}
}

func (r *Room) closeIfEmpty() {
	if r.onEmpty != nil {
		r.onEmpty(r.eventID, r)
	}
}
