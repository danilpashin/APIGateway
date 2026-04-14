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
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
