package eventusers

import (
	// Std

	// Intern

	// Extern
	"gorm.io/gorm"
)

type EventUser struct {
	gorm.Model

	UserID  uint `gorm:"not null"`
    EventID uint `gorm:"not null"`

	Role	string	`gorm:"not null;"`	// admin, member
}
