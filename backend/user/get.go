package user

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"gorm.io/gorm"
)

// registerGetUser

func registerGetUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "get-user",
        Method:         http.MethodGet,
        Path:           "/api/users/{name}",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleGetUser)
}

func (h *handler) handleGetUser(ctx context.Context, in *getUserInput) (*userOutput, error) {
    u, err := gorm.G[User](h.db).Where("name = ?", in.Name).First(ctx)
    if err != nil {
        return nil, err
    }
    return &userOutput{Body: u.ToSummaryDTO()}, nil
}

type getUserInput struct {
	Name string `path:"name" maxLength:"30" example:"Max" doc:"get user by name"`
}

// registerGetUsers

func registerGetUsers(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "get-users",
        Method:         http.MethodGet,
        Path:           "/api/users",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleGetUsers)
}

func (h *handler) handleGetUsers(ctx context.Context, in *getUsersInput) (*usersOutput, error) {
    offset := (in.Page - 1) * in.PageSize
    us, err := gorm.G[User](h.db).Limit(in.PageSize).Offset(offset).Find(ctx)
    if err != nil {
        return nil, err
    }
    
    userList := make([]UserSummaryDTO, 0, len(us))
    for _, u := range us {
        userList = append(userList, u.ToSummaryDTO())
    }

    out := usersOutput {
        Body: UserListSummaryDTO {
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
