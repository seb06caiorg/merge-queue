package models

import (
	"fmt"
	"strings"
	"time"
)

// User represents a user in the system.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // "admin", "user", "viewer"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

// UserFilter represents filtering options for users.
type UserFilter struct {
	Role     string `json:"role,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

// Validate checks if the user has valid data.
func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}
	if len(u.Username) < 3 || len(u.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !isValidEmail(u.Email) {
		return fmt.Errorf("invalid email format")
	}
	if !IsValidRole(u.Role) {
		return fmt.Errorf("invalid user role: %s", u.Role)
	}
	return nil
}

// IsValidRole checks if the role is valid.
func IsValidRole(role string) bool {
	validRoles := []string{"admin", "user", "viewer"}
	for _, v := range validRoles {
		if v == role {
			return true
		}
	}
	return false
}

// GetValidRoles returns all valid user roles.
func GetValidRoles() []string {
	return []string{"admin", "user", "viewer"}
}

// isValidEmail performs basic email validation.
func isValidEmail(email string) bool {
	// Basic email validation - in production, you'd want a proper regex or library.
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") &&
		strings.Contains(email, ".") &&
		len(email) > 5 &&
		len(email) < 255
}
