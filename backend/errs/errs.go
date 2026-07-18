package errs

import (
	// Std
	"errors"

	// External
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type ErrKind int

const (
	ErrNotFound ErrKind = iota
	ErrConflict
	ErrInvalidInput
	ErrInternal
	ErrCanNotRemoveAdmin
	ErrUserNotInEvent
)

func (k ErrKind) Error() string {
	switch k {
	case ErrNotFound:
		return "not found"
	case ErrConflict:
		return "conflict"
	case ErrInvalidInput:
		return "invalid input"
	case ErrInternal:
		return "internal error"
	case ErrCanNotRemoveAdmin:
		return "can not remove admin from event"
	case ErrUserNotInEvent:
		return "user is not registered for the event"
	default:
		return "unknown error"
	}
}

// Cama for Camaraderie
type CamaError struct {
	Kind    ErrKind
	Message string
}

func IsCamaError(err error) bool {
	var target *CamaError
	return errors.As(err, &target)
}

func NewCamaError(k ErrKind, msg string) CamaError {
	return CamaError{Kind: k, Message: msg}
}

func (e CamaError) Error() string {
	if e.Message == "" {
		return e.Kind.Error()
	}
	return e.Message
}

func CamaErrorString(err error) (string, bool) {
	var cErr *CamaError
	if errors.As(err, &cErr) {
		return cErr.Error(), true
	}

	return "", false
}

func (e CamaError) Is(target error) bool {
	k, ok := target.(ErrKind)
	return ok && e.Kind == k
}

func ErrorDB(err error) error {
	if err == nil {
		return nil
	}

	// postgres errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return NewCamaError(ErrConflict, "")
		case "23502", "23514": // can not be empty / check constraint violation
			return NewCamaError(ErrInvalidInput, "")
		}
	} else {
		return err
	}

	// gorm wrapped errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewCamaError(ErrNotFound, "")
	}

	return NewCamaError(ErrInternal, "")
}
