package user

import (
    // Std
	"time"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"gorm.io/gorm"
)

type handler struct {
    db *gorm.DB
}

func RegisterPublicApi(api huma.API, db *gorm.DB ) {
    db.AutoMigrate(&User{})

	h := handler{db: db}
    registerRegisterUser(api, h);
    registerLoginUser(api, h);
    registerLogoutUser(api, h);
}

func RegisterProtectedApi(api huma.API, db *gorm.DB ) {
    db.AutoMigrate(&User{})

	h := handler{db: db}
    registerGetUser(api, h);
    registerGetUsers(api, h);
    registerPatchUser(api, h);
}

type User struct {
    gorm.Model
    Name string     `gorm:"unique"`
	Email string    `gorm:"unique"`
    PasswordHash string
}

func (u *User) ToSummaryDTO() UserSummaryDTO {
    return UserSummaryDTO {
        ID:         u.ID,
        Name:       u.Name,
        Email:      u.Email,
        CreatedAt:  u.CreatedAt,
        UpdatedAt:  u.UpdatedAt,
    }
}

type userOutput struct {
    Body UserSummaryDTO
}

type usersOutput struct {
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
