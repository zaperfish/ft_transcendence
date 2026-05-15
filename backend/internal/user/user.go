package user

import (
    "ft_transcendence/backend/internal/db"
    "fmt"
	"context"
	// "net/http"

    "gorm.io/gorm"
    // "gorm.io/driver/postgres"

	"github.com/danielgtaylor/huma/v2"
	// "github.com/danielgtaylor/huma/v2/adapters/humachi"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

type handler struct {
    db *gorm.DB
}

func RegisterApi(api huma.API, db *gorm.DB) {
    h := &handler{db: db}
    huma.Get(api, "/api/users/{name}", h.HandleGet)
    huma.Post(api, "/api/users/new", h.HandleCreate)
}

type UserCreateInput struct {
    Body struct {
        Name string `json:"name" maxLength:"30" example:"Max" doc:"username"`
    }
}

type UserGetInput struct {
    Name string `path:"name" maxLength:"30" example:"1234" doc:"get user by id"`
}

type UserGetOutput struct {
    Body struct {
        Name string `json:"name" example:"Max" doc:"user creation confirmation"`
    }
}

type UserCreateOutput struct {
    Body struct {
        Message string `json:"message" example:"'Max' created successfully" doc:"user creation confirmation"`
    }
}

func (h *handler) HandleGet(ctx context.Context, input *UserGetInput) (*UserGetOutput, error) {
    resp := &UserGetOutput{}

    user, err := gorm.G[db.User](h.db).Where("name = ?", input.Name).First(ctx)
    if err != nil {
        return nil, err
    }

    // resp.Body.Id = user.Id
    resp.Body.Name = user.Name
    return resp, nil
}

func (h *handler) HandleCreate(ctx context.Context, input *UserCreateInput) (*UserCreateOutput, error) {
    resp := &UserCreateOutput{}
    fmt.Println("name:", input.Body.Name)
    err := gorm.G[db.User](h.db).Create(ctx, &db.User{Name: input.Body.Name})
    if err != nil {
        return nil, err
    }

    resp.Body.Message = fmt.Sprintf("'%s' created successfully", input.Body.Name)
    return resp, nil
}
