package domain

import "errors"

var (
	ErrProductsNotFound = errors.New("product(s) not found")
	ErrProductExist     = errors.New("product already exists")

	// required values
	ErrIDRequired = errors.New("id is required")

	// invalid values
	ErrInvalidID           = errors.New("id must be a number greater than 0")
	ErrInvalidName         = errors.New("name must be at least 2 characters and begin with Uppercase")
	ErrInvalidManufacturer = errors.New("manufacturer title must be at least 2 characters")
	ErrInvalidAmount       = errors.New("amount must be 0 or greater")
	ErrInvalidPrice        = errors.New("price must be greater than 0")
	ErrInvalidCategory     = errors.New("category must exist")
	ErrInvalidCursor       = errors.New("cursor is not a number")
	ErrInvalidLimit        = errors.New("limit is not a number")

	ErrNoUpdateData = errors.New("no updates provided")

	// query
	ErrQuery       = errors.New("error during the query")
	ErrCreateQuery = errors.New("invalid create query")
	ErrUpdateQuery = errors.New("invalid update query")
	ErrGetQuery    = errors.New("invalid get query")
	ErrListQuery   = errors.New("invalid list query")
	ErrDeleteQuery = errors.New("invalid delete query")

	// general
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidJSON  = errors.New("invalid JSON")
)

type ErrorResponse struct {
	Message string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}
