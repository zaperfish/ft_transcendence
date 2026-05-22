package main

import (
	// Std
	"log"
	"net/http"
	"os"

	// Internal
	"ft_transcendence/backend/app"
	// "ft_transcendence/backend/event"
	"ft_transcendence/backend/user"
	"ft_transcendence/backend/auth"

	// External
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("CONTAINER_RUNTIME") != "true" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Backend is local")
	} else {
		log.Println("Backend is in container")
	}

    app, err := app.Init()
	if err != nil {
		log.Fatal(err)
	}

    r := chi.NewRouter()
    
	config := huma.DefaultConfig("ft_transcendence api", "0.1.0")
	config.DocsRenderer = huma.DocsRendererScalar
	api := humachi.New(r, config)
	api.UseMiddleware(ChiMiddlewareToHuma(middleware.Logger))

    // Public Routes
	public := huma.NewGroup(api, "")
	auth.RegisterApi(public, app)

    // Protected Routes
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(ChiMiddlewareToHuma(jwtauth.Verifier(app.TokenAuth)))
	protected.UseMiddleware(ChiMiddlewareToHuma(jwtauth.Authenticator(app.TokenAuth)))
	user.RegisterApi(protected, app)

	startServer(r)
}

func ChiMiddlewareToHuma(chiMiddleware func(http.Handler) http.Handler) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		r, w := humachi.Unwrap(ctx)

		chiMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(humachi.NewContext(&huma.Operation{}, r, w))
		})).ServeHTTP(w, r)
	}
}

func startServer(r *chi.Mux) {
	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "4000"
	}

	log.Println("Listening on :" + port + "...")
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
