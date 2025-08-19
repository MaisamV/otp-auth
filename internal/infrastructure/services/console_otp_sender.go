package services

import (
	"context"
	"fmt"
	"log"

	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
)

// ConsoleOTPSender implements OTPSender interface for development/testing
// It prints OTP codes to the console instead of sending them via SMS
type ConsoleOTPSender struct {
	logger *log.Logger
}

// NewConsoleOTPSender creates a new console OTP sender
func NewConsoleOTPSender(logger *log.Logger) *ConsoleOTPSender {
	if logger == nil {
		logger = log.Default()
	}
	return &ConsoleOTPSender{
		logger: logger,
	}
}

// SendOTP prints the OTP to console (for development/testing)
func (s *ConsoleOTPSender) SendOTP(ctx context.Context, phoneNumber valueobjects.PhoneNumber, code string) error {
	if phoneNumber == "" {
		return errors.NewValidationError("Phone number is required", nil)
	}

	if code == "" {
		return errors.NewValidationError("OTP code is required", nil)
	}

	// Log the OTP to console
	message := fmt.Sprintf("[OTP SENDER] Sending OTP to %s: %s", string(phoneNumber), code)
	
	s.logger.Println(message)
	
	// Also print to stdout for visibility
	fmt.Printf("\n=== OTP NOTIFICATION ===\n")
	fmt.Printf("Phone: %s\n", string(phoneNumber))
	fmt.Printf("Code: %s\n", code)
	fmt.Printf("========================\n\n")

	return nil
}

// ValidatePhoneNumber validates if the phone number format is acceptable
func (s *ConsoleOTPSender) ValidatePhoneNumber(phoneNumber string) error {
	if phoneNumber == "" {
		return errors.NewValidationError("Phone number is required", nil)
	}

	// For console sender, we accept any non-empty phone number
	// In a real SMS service, you would validate the format more strictly
	if len(phoneNumber) < 10 {
		return errors.NewValidationError("Phone number must be at least 10 digits", nil)
	}

	return nil
}

// GetSenderInfo returns information about this OTP sender
func (s *ConsoleOTPSender) GetSenderInfo() map[string]interface{} {
	return map[string]interface{}{
		"type":        "console",
		"description": "Console OTP sender for development/testing",
		"enabled":     true,
	}
}