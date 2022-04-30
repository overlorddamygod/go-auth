package utils

import (
	"net/mail"

	"github.com/go-playground/validator/v10"
)

func IsEmailValid(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "This field must be at least " + fe.Param() + " characters"
	case "max":
		return "This field must be at most " + fe.Param() + " characters"
	}

	return fe.Error() // default error
}
