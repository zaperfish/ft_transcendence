package user

import (
	"time"

    "gorm.io/gorm"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func RegisterApi(api huma.API, db *gorm.DB) {
    db.AutoMigrate(&user{})

    h := dbHandler{db: db}
    registerCreateUser(api, h);
    registerGetUser(api, h);
    registerGetUsers(api, h);
    registerLoginUser(api, h);
}

type user struct {
    gorm.Model
    Name string     `gorm:"unique"`
    Password string
}

func (u *user) toResponseDTO() UserResponseDTO {
    return UserResponseDTO {
        ID:         u.ID,
        Name:       u.Name,
        CreatedAt:  u.CreatedAt,
        UpdatedAt:  u.UpdatedAt,
    }
}

type dbHandler struct {
    db *gorm.DB
}

type userOutput struct {
    Body UserResponseDTO
}

type usersOutput struct {
    Body UserListResponseDTO
}

type UserListResponseDTO struct {
    Data        []UserResponseDTO   `json:"data"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
	Total       int                 `json:"total"`
}

type UserResponseDTO struct {
    ID          uint        `json:"id" doc:"user ID"`
    Name        string      `json:"name" doc:"username"`
	CreatedAt   time.Time   `json:"created_at" doc:"user creation time"`
	UpdatedAt   time.Time   `json:"updated_at" doc:"user update time"`
}
