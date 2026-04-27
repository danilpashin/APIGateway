package domain

import (
	"time"
)

// ===== PRODUCT =====
type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Manufacturer string    `json:"manufacturer"`
	Price        int       `json:"price"`
	Amount       int       `json:"amount"`
	Status       bool      `json:"status"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
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
	CreatedAt    time.Time `json:"createdAt"`
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
	UpdatedAt    time.Time `json:"updatedAt"`
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
	NextCursor int    `json:"nextCursor"`
	HasMore    bool   `json:"hasMore"`
	Limit      uint64 `json:"limit"`
}
