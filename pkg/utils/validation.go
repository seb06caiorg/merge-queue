package utils

import (
	"fmt"
	"strings"
)

// ValidationUtils provides validation helper functions.
type ValidationUtils struct{}

// NewValidationUtils creates a new ValidationUtils instance.
func NewValidationUtils() *ValidationUtils {
	return &ValidationUtils{}
}

// IsValidEmail performs basic email validation.
func (vu *ValidationUtils) IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") &&
		strings.Contains(email, ".") &&
		len(email) > 5 &&
		len(email) < 255
}

// IsValidUsername checks if username meets basic requirements.
func (vu *ValidationUtils) IsValidUsername(username string) bool {
	username = strings.TrimSpace(username)
	return len(username) >= 3 && len(username) <= 50 && username != ""
}

// IsEmpty checks if a string is empty or only whitespace.
func (vu *ValidationUtils) IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Contains checks if a slice contains a string.
func (vu *ValidationUtils) Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ValidateRequired checks if required fields are present.
func (vu *ValidationUtils) ValidateRequired(fieldName, value string) error {
	if vu.IsEmpty(value) {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateLength checks if a string is within the specified length limits.
func (vu *ValidationUtils) ValidateLength(fieldName, value string, min, max int) error {
	length := len(strings.TrimSpace(value))
	if length < min {
		return fmt.Errorf("%s must be at least %d characters", fieldName, min)
	}
	if max > 0 && length > max {
		return fmt.Errorf("%s must be no more than %d characters", fieldName, max)
	}
	return nil
}

// ValidateOneOf checks if a value is one of the allowed values.
func (vu *ValidationUtils) ValidateOneOf(fieldName, value string, allowed []string) error {
	if !vu.Contains(allowed, value) {
		return fmt.Errorf("%s must be one of: %s", fieldName, strings.Join(allowed, ", "))
	}
	return nil
}

// SanitizeString removes leading/trailing whitespace and converts to lowercase.
func (vu *ValidationUtils) SanitizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ValidateTagList validates a list of tags.
func (vu *ValidationUtils) ValidateTagList(tags []string, maxTags int, maxTagLength int) error {
	if len(tags) > maxTags {
		return fmt.Errorf("maximum of %d tags allowed", maxTags)
	}

	for i, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return fmt.Errorf("tag %d is empty", i+1)
		}
		if len(tag) > maxTagLength {
			return fmt.Errorf("tag '%s' exceeds maximum length of %d characters", tag, maxTagLength)
		}
	}

	return nil
}
