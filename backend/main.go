package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
	)

	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	in_container_runtime := os.Getenv("CONTAINER_RUNTIME")
	if in_container_runtime != "true" {
		godotenv.Load()
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

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

	fmt.Println("Start listening...")
	err = http.ListenAndServe(":7772", r)
	if err != nil {
		log.Fatal(err)
	}
}
