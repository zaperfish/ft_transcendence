package main

import (
	// Std
	"log"
	"net/http"
	"os"

	// Internal
	"ft_transcendence/backend/db"
	"ft_transcendence/backend/event"
	"ft_transcendence/backend/user"
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/util"

	// External
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

    r := chi.NewRouter()
    
	config := huma.DefaultConfig("ft_transcendence api", "0.1.0")
	config.DocsRenderer = huma.DocsRendererScalar
	api := humachi.New(r, config)
	api.UseMiddleware(util.ChiMiddlewareToHuma(middleware.Logger))

    // Public Routes
	public := huma.NewGroup(api, "")
	auth.RegisterApi(public, db)

    // Protected Routes
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(auth.Verifier)
	protected.UseMiddleware(auth.Authenticator(api))
	user.RegisterApi(protected, db)
	event.RegisterEventsApi(protected, db)

	startServer(r)
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
