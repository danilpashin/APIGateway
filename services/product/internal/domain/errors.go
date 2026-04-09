package domain

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrForbidden       = errors.New("forbidden")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
