package middleware

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type PayloadValidator struct {
	Validator *validator.Validate
}

func CustomValidatorErrorMessage(err error) error {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErr {
			switch fieldErr.Tag() {
			case "required":
				return errors.New("Invalid request payload");
			}
		}
	}

	return nil;
}

func (pv *PayloadValidator) Validate(i interface{}) error {
	return pv.Validator.Struct(i);
}