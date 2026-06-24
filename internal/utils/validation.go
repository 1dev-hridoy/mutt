package utils

import (
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = "This field is required"
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = "Too short"
			case "max":
				errors[field] = "Too long"
			case "oneof":
				errors[field] = "Invalid value"
			default:
				errors[field] = "Invalid value"
			}
		}
		return errors
	}
	return map[string]string{"error": err.Error()}
}
