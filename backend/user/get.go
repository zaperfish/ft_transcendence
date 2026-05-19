package user

import (
	"net/http"
	"context"

    "gorm.io/gorm"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// registerGetUser

func registerGetUser(api huma.API, h dbHandler) {
    huma.Register(api, huma.Operation{
        OperationID:    "get-user",
        Method:         http.MethodGet,
        Path:           "/api/users/{name}",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleGetUser)
}

func (h *dbHandler) handleGetUser(ctx context.Context, in *getUserInput) (*userOutput, error) {
    u, err := gorm.G[user](h.db).Where("name = ?", in.Name).First(ctx)
    if err != nil {
        return nil, err
    }
    return &userOutput{Body: u.toResponseDTO()}, nil
}

type getUserInput struct {
    Name string `path:"name" maxLength:"30" example:"Max" doc:"get user by name"`
}

// registerGetUsers

func registerGetUsers(api huma.API, h dbHandler) {
    huma.Register(api, huma.Operation{
        OperationID:    "get-users",
        Method:         http.MethodGet,
        Path:           "/api/users",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleGetUsers)
}

func (h *dbHandler) handleGetUsers(ctx context.Context, in *getUsersInput) (*usersOutput, error) {
    offset := (in.Page - 1) * in.PageSize
    us, err := gorm.G[user](h.db).Limit(in.PageSize).Offset(offset).Find(ctx)
    if err != nil {
        return nil, err
    }
    
    userList := make([]userResponseDTO, 0, len(us))
    for _, u := range us {
        userList = append(userList, u.toResponseDTO())
    }

    out := usersOutput {
        Body: userListResponseDTO {
            Data:       userList,
            Page:       in.Page,
            PageSize:   in.PageSize,
            Total:      len(us),
        },
    }

    return &out, nil
}

type getUsersInput struct {
	Page     int `query:"page" minimum:"1" default:"1" doc:"Filter by page"`
	PageSize int `query:"page_size" minimum:"1" default:"10" doc:"Page size"`
}
