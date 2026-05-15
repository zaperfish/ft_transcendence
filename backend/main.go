package main

import (
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name string
}

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to DB:", dsn)

    db.AutoMigrate(&User{})

	return db, nil
}

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

	db, err := connectDB()
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

	huma.Get(api, "/api/postgres-version", h.HandlePostgresVersion)
	huma.Get(api, "/api/greeting/{name}", h.HandleGreeting)
    huma.Get(api, "/api/users/{name}", h.HandleUserGet)
    huma.Post(api, "/api/users/new", h.HandleCreateUser)

	log.Println("Listening on :7772...")
	err = http.ListenAndServe(":7772", r)
	if err != nil {
		log.Fatal(err)
	}
}

type Handler struct {
	db *gorm.DB
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

type UserCreateInput struct {
    Body struct {
        Name string `json:"name" maxLength:"30" example:"Max" doc:"username"`
    }
}

type UserGetInput struct {
    Name string `path:"name" maxLength:"30" example:"1234" doc:"get user by id"`
}

type UserGetOutput struct {
    Body struct {
        Name string `json:"name" example:"Max" doc:"user creation confirmation"`
    }
}

type UserCreateOutput struct {
    Body struct {
        Message string `json:"message" example:"'Max' created successfully" doc:"user creation confirmation"`
    }
}

func (h *Handler) HandleUserGet(ctx context.Context, input *UserGetInput) (*UserGetOutput, error) {
    resp := &UserGetOutput{}

    user, err := gorm.G[User](h.db).Where("name = ?", input.Name).First(ctx)
    if err != nil {
        return nil, err
    }

    // resp.Body.Id = user.Id
    resp.Body.Name = user.Name
    return resp, nil
}

func (h *Handler) HandleCreateUser(ctx context.Context, input *UserCreateInput) (*UserCreateOutput, error) {
    resp := &UserCreateOutput{}
    fmt.Println("name:", input.Body.Name)
    err := gorm.G[User](h.db).Create(ctx, &User{Name: input.Body.Name})
    if err != nil {
        return nil, err
    }

    resp.Body.Message = fmt.Sprintf("'%s' created successfully", input.Body.Name)
    return resp, nil
}
