package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	errMap := make(map[string]string, 0)

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()

		switch tag {
		case "required":
			errMap[field] = "this field is required"
		case "min":
			errMap[field] = "too short"
		case "max":
			errMap[field] = "too large"
		case "gt":
			errMap[field] = fmt.Sprint("must be greater than ", err.Param())
		case "gte":
			errMap[field] = fmt.Sprint("must be greater or equal ", err.Param())
		default:
			errMap[field] = "incorrect input"
		}
	}
	return errMap
}
