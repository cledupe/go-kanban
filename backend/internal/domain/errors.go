package domain

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict    = errors.New("conflict")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
