package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() (*gorm.DB, error) {
	dsn := "user=testuser password=securepass dbname=ft_transcendence host=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			var version string
			if err := db.Raw("SELECT version()").Scan(&version).Error; err != nil {
				http.Error(w, "query failed", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"version": version})
			// data := map[string]string{"message": "This message got server by the backend"}
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(data)
		})
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Healthcheck ok"))
		})
	})

	err = http.ListenAndServe(":7772", r)
	if err != nil {
		log.Fatal(err)
	}
}
