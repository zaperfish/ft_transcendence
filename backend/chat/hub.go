package chat

type Hub struct {
    rooms map[uint]*Room
}

func NewHub() *Hub {
    return &Hub{
        rooms: make(map[uint]*Room),
    }
}
