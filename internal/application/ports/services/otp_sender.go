package services

import (
	"context"

	"github.com/otp-auth/internal/domain/valueobjects"
)

// OTPSender defines the interface for sending OTP codes
type OTPSender interface {
	// SendOTP sends an OTP code to the specified phone number
	SendOTP(ctx context.Context, phoneNumber valueobjects.PhoneNumber, code string) error
}