package auth

import (
    // Std
	"context"
	// "fmt"
	"net/http"
	// "time"

    // Internal
	// "ft_transcendence/backend/app"
	// "ft_transcendence/backend/user"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	// "github.com/danielgtaylor/huma/v2/adapters/humachi"
    "github.com/go-chi/jwtauth/v5"
)

func MyMiddleware(ctx huma.Context, next func(huma.Context)) {
	next(ctx)
}

func Authenticator(ja *jwtauth.JWTAuth, api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		token, claims, err := jwtauth.FromContext(ctx.Context())
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid access token" , err)
			return
        }
		if token == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid access token")
			return
		}
		newCtx := context.WithValue(ctx.Context(), "claims", claims)
		next(huma.WithContext(ctx, newCtx))
	}
}
