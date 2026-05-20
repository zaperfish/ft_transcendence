package app

import (
	// Std
	"fmt"
	"log"
	"os"

	// External
    "github.com/go-chi/jwtauth/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const KEY = "secret"

type App struct {
    DB *gorm.DB
    TokenAuth *jwtauth.JWTAuth
}

func Init() (*App, error) {
    db, err := ConnectDB()
    if err != nil {
        return nil, err
    }

    a := &App {
        DB:         db,
        TokenAuth:  jwtauth.New("HS256", []byte(KEY), nil),
    }

    return a, nil
}

func ConnectDB() (*gorm.DB, error) {
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

	return db, nil
}
