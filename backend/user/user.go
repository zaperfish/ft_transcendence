package user

import (
    // Std
	"time"

    // External
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
	Name    string `gorm:"uniqueIndex:idx_name_active,where:deleted_at IS NULL"`
	Email   string `gorm:"uniqueIndex:idx_name_active,where:deleted_at IS NULL"`
    PasswordHash string
}

func (u *User) ToSummaryDTO() *UserSummaryDTO {
    return &UserSummaryDTO {
        ID:         u.ID,
        Name:       u.Name,
        Email:      u.Email,
        CreatedAt:  u.CreatedAt,
        UpdatedAt:  u.UpdatedAt,
    }
}

type UserOutput struct {
    Body UserSummaryDTO
}

type UsersOutput struct {
    Body UserListSummaryDTO
}

type UserListSummaryDTO struct {
    Data        []UserSummaryDTO   `json:"data"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
	Total       int                 `json:"total"`
}

type UserSummaryDTO struct {
    ID          uint        `json:"id" doc:"user ID"`
    Name        string      `json:"name" doc:"username"`
	Email 		string      `json:"email" doc:"email address"`
	CreatedAt   time.Time   `json:"created_at" doc:"user creation time"`
	UpdatedAt   time.Time   `json:"updated_at" doc:"user update time"`
}
