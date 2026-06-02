package main

import (
	// Std
	"log"
	"net/http"
	"os"

	// Internal
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/chat"
	"ft_transcendence/backend/db"
	"ft_transcendence/backend/event"
	"ft_transcendence/backend/middleware"
	"ft_transcendence/backend/user"

	// External
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
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

	err := auth.Init()
	if err != nil {
		log.Fatal(err)
	}

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
	err = db.AutoMigrate(&chat.Message{})
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)

	limiterStore := middleware.LimiterStore{
		IpLimiters:   make(map[string]*rate.Limiter),
		UserLimiters: make(map[string]*rate.Limiter),
	}

	r.Use(middleware.RateLimiterMiddleware(&limiterStore))
	r.Use(chiMiddleware.Logger)

	config := huma.DefaultConfig("ft_transcendence api", "0.1.0")
	config.DocsRenderer = huma.DocsRendererScalar
	config.CreateHooks = nil // disables schema injection into request json payloads
	api := humachi.New(r, config)

	// Public Routes
	public := huma.NewGroup(api, "")

	// Setup layers
	eventRepo := event.NewEventRepository(db)
	eventService := event.NewEventService(eventRepo)
	eventHandler := event.NewEventHandler(eventService)

	// Register routes
	// Register POST /events
	huma.Register(api, huma.Operation{
		OperationID:   "create-event",
		Method:        http.MethodPost,
		Path:          "/api/events",
		Summary:       "Create event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusCreated,
	}, eventHandler.CreateEvent)

	// Register PATCH /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "update-event",
		Method:        http.MethodPatch,
		Path:          "/api/events/{id}",
		Summary:       "Update event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, eventHandler.UpdateEvent)

	// Register DELETE /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "delete-event",
		Method:        http.MethodDelete,
		Path:          "/api/events/{id}",
		Summary:       "Delete event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, eventHandler.DeleteEvent)

	// Register GET /events/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "get-event",
		Method:        http.MethodGet,
		Path:          "/api/events/{id}",
		Summary:       "Get event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, eventHandler.GetEvent)

	// Register GET /events
	huma.Register(api, huma.Operation{
		OperationID:   "list-events",
		Method:        http.MethodGet,
		Path:          "/api/events",
		Summary:       "List events",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusOK,
	}, eventHandler.ListEvents)

	user.RegisterPublicRoutes(public, user.Handler{DB: db})

	// Protected Routes
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(auth.Verifier(api))
	protected.UseMiddleware(auth.Refresher(api))
	user.RegisterProtectedRoutes(protected, user.Handler{DB: db})
<<<<<<< Updated upstream
	chat.RegisterProtectedRoutes(protected, chat.NewHandler(db))
=======

	chatHandler := chat.NewHandler(db)
	chat.RegisterProtectedRoutes(protected, chatHandler)
	chat.RegisterWebSocketRoutes(r, chatHandler)
>>>>>>> Stashed changes

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
