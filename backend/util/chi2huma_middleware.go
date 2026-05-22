package util

import (
    // Std
	"net/http"

    // External
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
)

func ChiMiddlewareToHuma(chiMiddleware func(http.Handler) http.Handler) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		r, w := humachi.Unwrap(ctx)

		chiMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(humachi.NewContext(&huma.Operation{}, r, w))
		})).ServeHTTP(w, r)
	}
}
