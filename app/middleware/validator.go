package middleware

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type PayloadValidator struct {
	Validator *validator.Validate
	CustomValidatorErr func(err error) string
}

func CustomValidatorErrorMessage(err error) string {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErr {
			switch fieldErr.Tag() {
			case "required":
				return "Request not valid";
			}
		}
	}
	
	return err.Error();
}

func (cv *PayloadValidator) Validate(i interface{}) error {
    if err := cv.Validator.Struct(i); err != nil {
        return errors.New(cv.CustomValidatorErr(err));
    }
    return nil
}