package db

import (
	// Std
	"errors"
	"fmt"
	"log"
	"os"

	// External
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func PostgresError(err error) (error, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Println(pgErr)
		switch pgErr.Code {
		case "23505":
			return huma.Error409Conflict("already exists"), true
		case "23502":
			return huma.Error400BadRequest("can not be empty"), true
		}
	}
	return err, false
}
