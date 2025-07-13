package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func formatValidationErrors(errs validator.ValidationErrors) []ValidationError {
	var validationErrors []ValidationError
	for _, err := range errs {
		field := strings.ToLower(err.Field())
		switch err.Tag() {
		case "required":
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: "This field is required"})
		case "email":
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: "Invalid email format"})
		case "min":
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: fmt.Sprintf("Must be at least %s characters long", err.Param())})
		case "max":
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: fmt.Sprintf("Must be at most %s characters long", err.Param())})
		case "gte":
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: fmt.Sprintf("Must be greater than or equal to %s", err.Param())})
		case "lte":
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: fmt.Sprintf("Must be less than or equal to %s", err.Param())})
		case "oneof":
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: fmt.Sprintf("Must be one of: %s", err.Param())})
		default:
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: "Invalid value"})
		}
	}
	return validationErrors
}

func GetValidationErrors(err error) []ValidationError {
	if ve, ok := err.(validator.ValidationErrors); ok {
		return formatValidationErrors(ve)
	}
	return nil
}
