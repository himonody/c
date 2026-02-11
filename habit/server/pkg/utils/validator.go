package utils

import (
	"regexp"
)

// ValidateUsername validates username format (6-12 characters, alphanumeric and special characters)
func ValidateUsername(username string) bool {
	if len(username) < 6 || len(username) > 12 {
		return false
	}
	// Allow alphanumeric and special characters (letters, numbers, and common special chars)
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+$`, username)
	return matched
}

// ValidatePassword validates password format (6-12 characters, alphanumeric and special characters)
func ValidatePassword(password string) bool {
	if len(password) < 6 || len(password) > 12 {
		return false
	}
	// Allow alphanumeric and special characters
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+$`, password)
	return matched
}

// GetUsernameValidationError returns validation error message for username
func GetUsernameValidationError(username string) string {
	if len(username) < 6 {
		return "username must be at least 6 characters"
	}
	if len(username) > 12 {
		return "username must not exceed 12 characters"
	}
	if !ValidateUsername(username) {
		return "username can only contain letters, numbers, and special characters"
	}
	return ""
}

// GetPasswordValidationError returns validation error message for password
func GetPasswordValidationError(password string) string {
	if len(password) < 6 {
		return "password must be at least 6 characters"
	}
	if len(password) > 12 {
		return "password must not exceed 12 characters"
	}
	if !ValidatePassword(password) {
		return "password can only contain letters, numbers, and special characters"
	}
	return ""
}
