package domain

import "errors"

var (
	ErrProductsNotFound = errors.New("products not found")
	ErrProductExist     = errors.New("product already exists")

	ErrNameRequired         = errors.New("name is required")
	ErrManufacturerRequired = errors.New("manufacturer is required")
	ErrPriceRequired        = errors.New("price is required")
	ErrAmountRequired       = errors.New("amount is required")
	ErrCategoryRequired     = errors.New("category is required")

	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
)

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}
