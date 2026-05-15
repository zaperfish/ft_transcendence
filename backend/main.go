package main

import (
    "ft_transcendence/backend/user"
    "ft_transcendence/backend/db"
	"context"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	// "github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func main() {
	if os.Getenv("LOCAL_DEV") == "true" {
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
	r.Use(middleware.Logger)

	config := huma.DefaultConfig("ft_transcendence api", "0.1.0")
	config.DocsRenderer = huma.DocsRendererScalar
	// config.SchemasPath = ""
	// config.OpenAPIPath = "/api/openapi.json"
	// config.DocsPath = ""

	api := humachi.New(r, config)

	h := &Handler{db: db}

    user.RegisterApi(api, db)
	huma.Get(api, "/api/postgres-version", h.HandlePostgresVersion)
	huma.Get(api, "/api/greeting/{name}", h.HandleGreeting)

	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "4001"
	}

	log.Println("Listening on :" + port + "...")
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}

type PostgresVersionOutput struct {
	Body struct {
		Version string `json:"version"`
	}
}

func (h *Handler) HandlePostgresVersion(ctx context.Context, input *struct{}) (*PostgresVersionOutput, error) {
	var version string
	err := h.db.Raw("SELECT version()").Scan(&version).Error
	if err != nil {
		return nil, huma.Error500InternalServerError("query failed")
	}
	resp := &PostgresVersionOutput{}
	resp.Body.Version = version
	return resp, nil
}

type GreetingInput struct {
	Name string `path:"name" maxLength:"30" example:"Max" doc:"Name to greet"`
}

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func (h *Handler) HandleGreeting(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
	resp := &GreetingOutput{}
	resp.Body.Message = fmt.Sprintf("Hello %s", input.Name)
	return resp, nil
}
