package domain

import "errors"

var (
	// Business errors
	ErrUserNotFound         = errors.New("user not found")
	ErrRoleNotFound         = errors.New("role not found")
	ErrUserExist            = errors.New("user already exists")
	ErrWrongEmailOrPassword = errors.New("invalid email or password")

	// Validation errors
	ErrIDRequired       = errors.New("id is required")
	ErrUsernameRequired = errors.New("username is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrInvalidID        = errors.New("id must be a positive number")
	ErrInvalidUsername  = errors.New("username must be at least 2 characters")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrInvalidPassword  = errors.New("password must be at least 8 characters")
	ErrInvalidCursor    = errors.New("invalid cursor value")
	ErrInvalidLimit     = errors.New("limit must be between 1 and 50")
	ErrInvalidJSON      = errors.New("invalid json format")
	ErrNoInsertData     = errors.New("required fields are missing")
	ErrNoUpdateData     = errors.New("no update values")
)

type ErrorResponse struct {
	Message string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}
