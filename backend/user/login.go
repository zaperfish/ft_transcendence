package user

import (
	"net/http"
	"fmt"
	"context"

    "gorm.io/gorm"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

type userLoginDTO struct {
    Name string     `json:"name"`
    Password string `json:"password"`
}

func registerLoginUser(api huma.API, h dbHandler) {
    huma.Register(api, huma.Operation{
        OperationID:    "login-user",
        Method:         http.MethodPost,
        Path:           "/api/auth/login",
        DefaultStatus:  http.StatusOK,
        Tags:           []string{"Users"},
    }, h.handleLoginUser)
}

type loginUserInput struct {
    Body userLoginDTO
}

// type userLoginResponseDTO struct {
//     responseDTO userResponseDTO
//     accessToken string              `json:"access_token"`
//     tokenType string                `json:"token_type"`
// }
//
// type loginUserOutput struct {
//     Body userLoginResponseDTO
// }

type UserLoginResponseDTO struct {
    AccessToken string              `json:"access_token"`
    TokenType string                `json:"token_type"`
}

type loginUserOutput struct {
    Body UserLoginResponseDTO
}

func (h *dbHandler) handleLoginUser(ctx context.Context, in *loginUserInput) (*loginUserOutput, error) {
    u, err := gorm.G[user](h.db).Where("name = ?", in.Body.Name).First(ctx)
    if err != nil {
        return nil, err
    }

    if u.Password != in.Body.Password {
        return nil, gorm.ErrRecordNotFound
    }
    // create jwt
    
    out := loginUserOutput {
        Body: UserLoginResponseDTO {
            AccessToken:    "testToken",
            TokenType:      "Bearer",
        },
    }
    fmt.Printf("OUT: %v\n", out)
    return &out, nil
}
