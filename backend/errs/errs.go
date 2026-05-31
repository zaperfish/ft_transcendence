package errs

import (
	// Std
	"errors"

	// External
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal error")
)

func ErrorDB(err error) error {
	if err == nil {
		return nil
	}

	// postgres errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrConflict
		case "23502", "23514": // can not be empty / check constraint violation
			return ErrInvalidInput
		}
	}

	// gorm wrapped errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	return ErrInternal
}
