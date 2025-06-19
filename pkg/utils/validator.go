package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMessages []string
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			switch fieldError.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", field))
			case "email":
				errorMessages = append(errorMessages, "invalid email format")
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long", field, fieldError.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at most %s characters long", field, fieldError.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("invalid value for %s", field))
			}
		}
		if len(errorMessages) > 0 {
			msg := strings.Join(errorMessages, ", ")
			return strings.ToUpper(string(msg[0])) + msg[1:] + "."
		}
	}
	return "Invalid input provided."
}
