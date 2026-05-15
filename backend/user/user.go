package user

import (
	"net/http"
    "ft_transcendence/backend/db"
	"context"
    "gorm.io/gorm"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

type handler struct {
    db *gorm.DB
}

func RegisterApi(api huma.API, database *gorm.DB) {
    registerGet(api, database)
    registerCreate(api, database)
}

func registerCreate(api huma.API, database *gorm.DB) {
    type input struct {
        Body struct {
            Name string `json:"name" maxLength:"30" example:"Max" doc:"username"`
        }
    }
    huma.Register(api, huma.Operation{
        OperationID:    "create-user",
        Method:         http.MethodPost,
        Path:           "/api/users/new",
        Summary:        "Create a new user",
        Description:    "Create a new user with all the parameters",
        DefaultStatus:  201,
        Tags:           []string{"Users"},
    }, func(ctx context.Context, in *input) (*struct{}, error) {
        err := gorm.G[db.User](database).Create(ctx, &db.User{Name: in.Body.Name})
        if err != nil {
            return nil, err
        }
        return nil, nil
    })
}

func registerGet(api huma.API, database *gorm.DB) {
    type input struct {
        Name string `path:"name" maxLength:"30" example:"1234" doc:"get user by id"`
    }
    type output struct {
        Body struct {
            Name string `json:"name" example:"Max" doc:"user creation confirmation"`
        }
    }
    huma.Register(api, huma.Operation{
        OperationID:    "get-user",
        Method:         http.MethodGet,
        Path:           "/api/users/{name}",
        Summary:        "Query user information",
        Description:    "Get all the information about a user",
        DefaultStatus:  200,
        Tags:           []string{"Users"},
    }, func(ctx context.Context, in *input) (*output, error) {
        resp := &output{}
        user, err := gorm.G[db.User](database).Where("name = ?", in.Name).First(ctx)
        if err != nil {
            return nil, err
        }
        resp.Body.Name = user.Name
        return resp, nil
    })
}
