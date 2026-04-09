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
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Price        int    `json:"price"`
	Amount       int    `json:"amount"`
	Category     string `json:"category"`
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
	Status       bool    `json:"status"`
	Category     string  `json:"category"`
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
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Price        int    `json:"price"`
	Category     string `json:"category"`
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
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Price        int    `json:"price"`
}

type ListProductsResponse struct {
	Products []*Product
	Total    int
}
