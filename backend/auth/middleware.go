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

// verifies existence and validity of token
// saves the verified token in context
func Verifier(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// get jwt cookie
		tokenCookie, err := huma.ReadCookie(ctx, "jwt")
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, err.Error())
			return
		}
		// verify and potentially extract token
		token, err := jwtauth.VerifyToken(tokenAuth, tokenCookie.Value)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, err.Error())
			return
		}
		// save token in context
		newCtx := jwtauth.NewContext(ctx.Context(), token, nil)
		next(huma.WithContext(ctx, newCtx))
	}
}

// refreshes token and saves claims in context
func Refresher(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		 // the below line depends on the Verifier having saved the token to the context
		_, claims, err := jwtauth.FromContext(ctx.Context())
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, err.Error())
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
