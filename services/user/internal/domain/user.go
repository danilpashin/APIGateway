package domain

import "time"

// ===== USER =====
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ===== CREATE =====
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type CreateUserResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// ===== UPDATE =====
type UpdateUserRequest struct {
	Username    *string `json:"username" validate:"omitnil,min=3,max=30"`
	Email       *string `json:"email" validate:"omitnil,email,max=254"`
	Password    *string `json:"oldPassword" validate:"omitnil,min=8,max=72"`
	NewPassword *string `json:"newPassword" validate:"omitnil,min=8,max=72"`
}

type UpdateUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ===== GET =====
type GetUserRequest struct{}

type GetUserResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

// ===== LIST =====
type ListUsersRequest struct{}

type ListUsersResponse struct {
	Users      []*User
	Pagination *Pagination
}

type Pagination struct {
	NextCursor int    `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
	Limit      uint64 `json:"limit"`
}

// ===== DELETE =====
type DeleteUserRequest struct{}

type DeleteUserResponse struct{}
