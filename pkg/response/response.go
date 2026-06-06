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
			errMap[field] = fmt.Sprintf("must be at least %s characters", err.Param())
		case "max":
			errMap[field] = fmt.Sprintf("must be at most %s characters", err.Param())
		case "gt":
			errMap[field] = fmt.Sprintf("must be greater than %s", err.Param())
		case "gte":
			errMap[field] = fmt.Sprintf("must be greater or equal %s", err.Param())
		case "email":
			errMap[field] = "invalid email format"
		default:
			errMap[field] = "incorrect input"
		}
	}
	return errMap
}
