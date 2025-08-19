package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhoneNumber validates a phone number format
func ValidatePhoneNumber(phone string) bool {
	// Remove any whitespace
	phone = strings.TrimSpace(phone)
	
	if phone == "" {
		return false
	}

	// Check international format: +<country_code><phone_number>
	if strings.HasPrefix(phone, "+") {
		internationalPattern := regexp.MustCompile(`^\+\d{1,3}9\d{9}$`)
		return internationalPattern.MatchString(phone)
	}

	// Check local format: 0<phone_number>
	if strings.HasPrefix(phone, "0") {
		localPattern := regexp.MustCompile(`^09\d{9}$`)
		return localPattern.MatchString(phone)
	}

	return false
}

// ValidatePassword validates password strength
func ValidatePassword(password string) bool {
	// At least 8 characters
	if len(password) < 8 {
		return false
	}
	
	// Contains at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Contains at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Contains at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	// Contains at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	
	return hasUpper && hasLower && hasDigit && hasSpecial
}

// SanitizeString removes potentially harmful characters from a string
func SanitizeString(input string) string {
	// Remove null bytes and control characters
	sanitized := regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
	// Trim whitespace
	return strings.TrimSpace(sanitized)
}

// ValidateUUID validates a UUID format
func ValidateUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}