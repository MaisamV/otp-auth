package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

// PhoneNumber represents a validated phone number
type PhoneNumber string

// Phone number validation patterns
var (
	// Pattern for +<country_code><phone_number> format
	internationalPattern = regexp.MustCompile(`^\+\d{1,3}9\d{9}$`)
	// Pattern for 0<phone_number> format (Iranian format)
	localPattern = regexp.MustCompile(`^09\d{9}$`)
)

// NewPhoneNumber creates and validates a new phone number
func NewPhoneNumber(phone string) (PhoneNumber, error) {
	// Remove any whitespace
	phone = strings.TrimSpace(phone)
	
	if phone == "" {
		return "", errors.New("phone number cannot be empty")
	}

	// Validate the phone number format
	if err := validatePhoneNumber(phone); err != nil {
		return "", err
	}

	// Normalize local numbers to international format
	if strings.HasPrefix(phone, "0") {
		// Convert 09xxxxxxxxx to +989xxxxxxxxx
		phone = "+98" + phone[1:]
	}

	return PhoneNumber(phone), nil
}

// validatePhoneNumber validates the phone number format
func validatePhoneNumber(phone string) error {
	// Check international format: +<country_code><phone_number>
	if strings.HasPrefix(phone, "+") {
		if !internationalPattern.MatchString(phone) {
			return errors.New("invalid international phone number format. Must be +<country_code><phone_number> where phone_number is 10 digits starting with 9")
		}
		return nil
	}

	// Check local format: 0<phone_number>
	if strings.HasPrefix(phone, "0") {
		if !localPattern.MatchString(phone) {
			return errors.New("invalid local phone number format. Must be 0<phone_number> where phone_number is 10 digits starting with 9")
		}
		return nil
	}

	return errors.New("phone number must start with + (international) or 0 (local)")
}

// String returns the string representation of the phone number in international format
func (p PhoneNumber) String() string {
	// Always return international format
	phone := string(p)
	if strings.HasPrefix(phone, "0") {
		// Convert 09xxxxxxxxx to +989xxxxxxxxx
		return "+98" + phone[1:]
	}
	return phone
}

// IsValid checks if the phone number is valid
func (p PhoneNumber) IsValid() bool {
	return validatePhoneNumber(string(p)) == nil
}

// ToInternational converts local format to international format (assumes Iran +98)
func (p PhoneNumber) ToInternational() PhoneNumber {
	phone := string(p)
	if strings.HasPrefix(phone, "0") {
		// Convert 09xxxxxxxxx to +989xxxxxxxxx
		return PhoneNumber("+98" + phone[1:])
	}
	return p
}

// ToLocal converts international format to local format (for Iranian numbers)
func (p PhoneNumber) ToLocal() PhoneNumber {
	phone := string(p)
	if strings.HasPrefix(phone, "+98") {
		// Convert +989xxxxxxxxx to 09xxxxxxxxx
		return PhoneNumber("0" + phone[3:])
	}
	return p
}