package domain

import (
	"time"
)

// ===== PRODUCT =====
type Product struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Manufacturer string    `json:"manufacturer" db:"manufacturer"`
	Price        int       `json:"price" db:"price"`
	Amount       int       `json:"amount" db:"amount"`
	Status       bool      `json:"status" db:"category"`
	Category     string    `json:"category" db:"category"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ===== CREATE =====
type CreateProductRequest struct {
	Name         string `json:"name" validate:"required,min=5,max=150"`
	Manufacturer string `json:"manufacturer" validate:"required,min=2,max=50"`
	Price        int    `json:"price" validate:"required,gt=0"`
	Amount       int    `json:"amount" validate:"required,gte=0"`
	Status       bool   `json:"status,omitempty"`
	Category     string `json:"category" validate:"required,min=5,max=100"`
}

type CreateProductResponse struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Manufacturer string    `json:"manufacturer"`
	Price        int       `json:"price"`
	Amount       int       `json:"amount"`
	Status       bool      `json:"status"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"created_at"`
}

// ===== UPDATE =====
type UpdateProductRequest struct {
	Name         *string `json:"name"`
	Manufacturer *string `json:"manufacturer"`
	Price        *int    `json:"price"`
	Amount       *int    `json:"amount"`
	Status       *bool   `json:"status"`
	Category     *string `json:"category"`
}

type UpdateProductResponse struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Manufacturer string    `json:"manufacturer"`
	Price        int       `json:"price"`
	Amount       int       `json:"amount"`
	Status       bool      `json:"status"`
	Category     string    `json:"category"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ===== GET =====
type GetProductRequest struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Price        int    `json:"price,omitempty"`
	Category     string `json:"category,omitempty"`
}

type GetProductResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Price        int    `json:"price"`
	Amount       int    `json:"amount"`
	Status       bool   `json:"status"`
	Category     string `json:"category"`
}

// ===== LIST =====
type ListProductsRequest struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Price        int    `json:"price,omitempty"`
}

type ListProductsResponse struct {
	Products         []*Product  `json:"data"`
	PaginationParams *Pagination `json:"pagination"`
	Total            int         `json:"total"`
}

type Pagination struct {
	NextCursor int    `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
	Limit      uint64 `json:"limit"`
}
