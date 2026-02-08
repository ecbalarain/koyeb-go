package validator

import (
	"html"
	"strings"
)

// Sanitize removes potentially dangerous characters and HTML tags from a string.
// This helps prevent XSS attacks by escaping HTML entities.
func Sanitize(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)
	
	// Escape HTML entities to prevent XSS
	sanitized := html.EscapeString(trimmed)
	
	return sanitized
}

// SanitizeMultiple sanitizes multiple strings at once.
func SanitizeMultiple(inputs ...string) []string {
	result := make([]string, len(inputs))
	for i, input := range inputs {
		result[i] = Sanitize(input)
	}
	return result
}

// ValidateAndSanitize validates that a string is not empty after sanitization.
func ValidateAndSanitize(input string, fieldName string) (string, error) {
	sanitized := Sanitize(input)
	if sanitized == "" {
		return "", &ValidationError{
			Field:   fieldName,
			Message: fieldName + " cannot be empty",
		}
	}
	return sanitized, nil
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
