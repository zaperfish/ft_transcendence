package event

import (
	// Std
	"time"
	// Intern
	// Extern
)

type Event struct {
	ID              uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Title           string
	Description     string
	StartTime       time.Time
	Duration        int
	LocationName    string
	LocationAddress string
	MaxCapacity     uint
	NumRegistered   uint
	ImagePath  		string
}
