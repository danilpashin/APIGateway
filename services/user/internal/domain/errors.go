package domain

import "errors"

var (
	// Business errors
	ErrUserNotFound         = errors.New("user not found")
	ErrRoleNotFound         = errors.New("role not found")
	ErrUserExist            = errors.New("user already exists")
	ErrWrongEmailOrPassword = errors.New("invalid email or password")

	// Validation errors
	ErrIDRequired      = errors.New("id is required")
	ErrInvalidID       = errors.New("id must be a positive number")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
	ErrInvalidCursor   = errors.New("invalid cursor value")
	ErrInvalidLimit    = errors.New("limit must be between 1 and 100")
	ErrInvalidJSON     = errors.New("invalid json format")
	ErrNoInsertData    = errors.New("required fields are missing")
)

type ErrorResponse struct {
	Message string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}
