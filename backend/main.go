package main

import (
	// Std
	"log"
	"net/http"
	"os"

	// Internal
	"ft_transcendence/backend/apikey"
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/chat"
	"ft_transcendence/backend/db"
	"ft_transcendence/backend/devtools"
	"ft_transcendence/backend/event"
	"ft_transcendence/backend/eventusers"
	"ft_transcendence/backend/health"
	"ft_transcendence/backend/me"
	"ft_transcendence/backend/middleware"
	"ft_transcendence/backend/user"

	// External
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
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

	db := initDB()
	r := chi.NewRouter()
	initApi(r, db)

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

func initDB() *gorm.DB {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&event.GormEventModel{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&eventusers.EventUser{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&chat.Message{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func initApi(r *chi.Mux, db *gorm.DB) {
	r.Use(chiMiddleware.Logger)

	err := auth.Init()
	if err != nil {
		log.Fatal(err)
	}

	limiterStore := middleware.LimiterStore{
		IpLimiters:   make(map[string]*rate.Limiter),
		UserLimiters: make(map[string]*rate.Limiter),
	}

	r.Use(middleware.RateLimiterMiddleware(&limiterStore))

	config := huma.DefaultConfig("ft_transcendence api", "0.1.0")
	config.DocsRenderer = huma.DocsRendererScalar
	config.CreateHooks = nil // disables schema injection into request json payloads
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"AdminPassword": {
			Type:        "http",
			Scheme:      "bearer",
			Description: "Enter admin password as: Bearer <password>",
		},
		"ApiKey": {
			Type:        "http",
			Scheme:      "bearer",
			Description: "Enter api key as: Bearer <key>",
		},
		"SessionToken": {
			Type:        "http",
			Scheme:      "bearer",
			Description: "JWT",
		},
	}

	api := humachi.New(r, config)

	// Generate and return the Prometheus metrics payload at /metrics
	r.Handle("/metrics", promhttp.Handler())
	devtools.RegisterRoutes(r)

	apikey.RegisterRoutes(api, db)
	chatHandler := chat.NewHandler(db)
	event.RegisterRoutes(api, db, chatHandler.Hub)
	health.RegisterRoutes(api)

	userHandler := user.NewHandler(db)
	meHandler := me.NewHandler(db, chatHandler.Hub)
	// Public Routes
	public := huma.NewGroup(api, "")

	user.RegisterPublicRoutes(public, userHandler)

	// Protected Routes
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(auth.Verifier(api))
	protected.UseMiddleware(auth.Refresher(api))
	user.RegisterProtectedRoutes(protected, userHandler)
	me.RegisterRoutes(protected, meHandler)

	chat.RegisterProtectedRoutes(protected, chatHandler)
	chat.RegisterWebSocketRoutes(r, chatHandler)
}
