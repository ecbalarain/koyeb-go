package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// Validate validates a struct and returns a user-friendly error message.
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("validation failed: %w", err)
	}

	return formatValidationErrors(validationErrors)
}

// formatValidationErrors converts validator errors into a user-friendly message.
func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string

	for _, err := range errs {
		msg := formatFieldError(err)
		messages = append(messages, msg)
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

// formatFieldError formats a single field validation error.
func formatFieldError(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	default:
		return fmt.Sprintf("%s failed validation (%s)", field, err.Tag())
	}
}

// ValidOrderStatus checks if the given status is valid for orders.
func ValidOrderStatus(status string) bool {
	validStatuses := map[string]bool{
		"new":       true,
		"confirmed": true,
		"shipped":   true,
		"canceled":  true,
	}
	return validStatuses[status]
}
