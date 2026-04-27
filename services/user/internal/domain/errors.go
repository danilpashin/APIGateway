package domain

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrRoleNotFound    = errors.New("role not found")
	ErrUserExist       = errors.New("user already exist")
	ErrIDRequired      = errors.New("id is required")
	ErrWrongPassword   = errors.New("wrong password")
	ErrNoInsertData    = errors.New("insert data is empty or one of the fields is not filled")
	ErrInvalidID       = errors.New("id must be a number greater than 0")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
	ErrInvalidJSON     = errors.New("invalid JSON")
)
