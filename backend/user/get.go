package user

import (
	"net/http"
	"context"

    "fmt"
    // "github.com/go-chi/jwtauth/v5"

    "gorm.io/gorm"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// registerGetUser

func registerGetUser(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "get-user",
        Method:         http.MethodGet,
        Path:           "/users/{name}",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleGetUser)
}

func (h *handler) handleGetUser(ctx context.Context, in *getUserInput) (*userOutput, error) {

	claims := ctx.Value("claims")

	fmt.Println("claims:\t", claims)

    u, err := gorm.G[User](h.db).Where("name = ?", in.Name).First(ctx)
    if err != nil {
        return nil, err
    }
    return &userOutput{Body: u.ToResponseDTO()}, nil
}

type getUserInput struct {
    Name string `path:"name" maxLength:"30" example:"Max" doc:"get user by name"`
}

// registerGetUsers

func registerGetUsers(api huma.API, h handler) {
    huma.Register(api, huma.Operation{
        OperationID:    "get-users",
        Method:         http.MethodGet,
        Path:           "/users",
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
    
    userList := make([]UserResponseDTO, 0, len(us))
    for _, u := range us {
        userList = append(userList, u.ToResponseDTO())
    }

    out := usersOutput {
        Body: UserListResponseDTO {
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
