package middleware

import (
	"github.com/go-playground/validator/v10"
)

type PayloadValidator struct {
	Validator *validator.Validate
}

func (cv *PayloadValidator) Validate(i interface{}) error {
    return cv.Validator.Struct(i);
}