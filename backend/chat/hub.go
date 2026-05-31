package chat

import (
	// Std
	"sync"
)

type Hub struct {
	mu sync.Mutex
	// rooms stores active chat rooms by event ID
	rooms map[uint]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[uint]*Room),
	}
}

func (h *Hub) GetOrCreateRoom(eventID uint) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.rooms[eventID]
	if ok {
		return room
	}

	room = newRoom(eventID, h.removeRoom)
	h.rooms[eventID] = room
	go room.run()

	return room
}

func (h *Hub) JoinRoom(eventID uint, client *Client) *Room {
	for {
		room := h.GetOrCreateRoom(eventID)
		if room.Join(client) {
			return room
		}
	}
}

func (h *Hub) removeRoom(eventID uint, room *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[eventID] == room {
		delete(h.rooms, eventID)
	}
}
