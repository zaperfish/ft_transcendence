package chat

type Hub struct {
	// rooms stores active chat rooms by event ID
	rooms map[uint]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[uint]*Room),
	}
}
