package auth

import (
    // Std
	"fmt"
	"net/http"

    // Internal
	// "ft_transcendence/backend/app"
	// "ft_transcendence/backend/user"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
    // "github.com/go-chi/jwtauth/v5"
    "github.com/go-chi/jwtauth/v5"
)

func MyMiddleware(ctx huma.Context, next func(huma.Context)) {
	next(ctx)
}

func Authenticator(ja *jwtauth.JWTAuth, api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		token, claims, err := jwtauth.FromContext(ctx.Context())
		fmt.Println("token:\t", token)
		fmt.Println("claims:\t", claims)
		fmt.Println("err:\t", err)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid access token" , err)
			return
        }
		if token == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid access token")
			return
		}
		next(huma.WithValue(ctx, "claims", claims))
	}
}
//
// func Authenticator(ja *JWTAuth) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		hfn := func(w http.ResponseWriter, r *http.Request) {
// 			token, _, err := FromContext(r.Context())
//
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusUnauthorized)
// 				return
// 			}
//
// 			if token == nil {
// 				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 				return
// 			}
//
// 			// Token is authenticated, pass it through
// 			next.ServeHTTP(w, r)
// 		}
// 		return http.HandlerFunc(hfn)
// 	}
// }
