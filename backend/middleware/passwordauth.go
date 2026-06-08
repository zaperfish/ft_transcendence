package middleware

import (
	// Std
	"crypto/subtle"
	"net/http"
	"strings"

	// External
	"github.com/danielgtaylor/huma/v2"
)

func PasswordVerifier(api huma.API, expectedPassword string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		auth := ctx.Header("Authorization")
		if auth == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid Authorization format")
			return
		}

		provided := strings.TrimPrefix(auth, "Bearer ")

		// constant-time compare (prevents timing attacks)
		if subtle.ConstantTimeCompare([]byte(provided), []byte(expectedPassword)) != 1 {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid password")
			return
		}

		next(ctx)
	}
}
