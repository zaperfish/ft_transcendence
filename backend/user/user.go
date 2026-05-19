package user

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"gorm.io/gorm"
)

func RegisterApi(api huma.API, db *gorm.DB) {
	db.AutoMigrate(&User{})

	h := dbHandler{db: db}
	registerCreateUser(api, h)
	registerGetUser(api, h)
	registerGetUsers(api, h)
}

type user struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func (u *User) toResponseDTO() userResponseDTO {
	return userResponseDTO{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type dbHandler struct {
	db *gorm.DB
}

type userOutput struct {
	Body userResponseDTO
}

type usersOutput struct {
	Body userListResponseDTO
}

type userListResponseDTO struct {
	Data     []userResponseDTO `json:"data"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int               `json:"total"`
}

type userResponseDTO struct {
	ID        uint      `json:"id" doc:"user ID"`
	Name      string    `json:"name" doc:"username"`
	CreatedAt time.Time `json:"created_at" doc:"user creation time"`
	UpdatedAt time.Time `json:"updated_at" doc:"user update time"`
}
