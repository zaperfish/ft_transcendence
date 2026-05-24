package auth

import (
    // Std
	"context"
	"net/http"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
    "github.com/go-chi/jwtauth/v5"
)

func Authenticator(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
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

		cookie, err := makeJWTCookie(claims["sub"].(string))
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "error")
			return
		}

		_, w := humachi.Unwrap(ctx)
		http.SetCookie(w, &cookie)
		newCtx := context.WithValue(ctx.Context(), "claims", claims)
		next(huma.WithContext(ctx, newCtx))
	}
}

func Verifier(ctx huma.Context, next func(huma.Context)) {
	chiMiddlewareToHuma(jwtauth.Verifier(tokenAuth))(ctx, next)
}

func chiMiddlewareToHuma(chiMiddleware func(http.Handler) http.Handler) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		r, w := humachi.Unwrap(ctx)
		chiMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(humachi.NewContext(&huma.Operation{}, r, w))
		})).ServeHTTP(w, r)
	}
}
