package validator

import "github.com/go-playground/validator/v10"

var v = validator.New()

func New(i interface{}) error {
	return v.Struct(i)
}
