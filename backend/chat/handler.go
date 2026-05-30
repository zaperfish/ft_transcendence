package chat

import "gorm.io/gorm"

type Handler struct {
    DB  *gorm.DB
    Hub *Hub
}

func NewHandler(db *gorm.DB) Handler {
    return Handler{
        DB:  db,
        Hub: NewHub(),
    }
}
