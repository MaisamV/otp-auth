package repositories

import (
	"context"
	"time"

	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
)

// OTPReader defines read operations for OTPs
type OTPReader interface {
	// Get retrieves an OTP by phone number
	Get(ctx context.Context, phoneNumber valueobjects.PhoneNumber) (*entities.OTP, error)
	
	// Exists checks if an OTP exists for the given phone number
	Exists(ctx context.Context, phoneNumber valueobjects.PhoneNumber) (bool, error)
}

// OTPWriter defines write operations for OTPs
type OTPWriter interface {
	// Store stores an OTP with TTL
	Store(ctx context.Context, otp *entities.OTP, ttl time.Duration) error
	
	// Delete deletes an OTP by phone number
	Delete(ctx context.Context, phoneNumber valueobjects.PhoneNumber) error
}

// OTPRepository combines read and write operations
type OTPRepository interface {
	OTPReader
	OTPWriter
}

// RateLimiter defines rate limiting operations
type RateLimiter interface {
	// CheckAndIncrement checks the current count and increments if under limit
	CheckAndIncrement(ctx context.Context, key string, limit int, window time.Duration) (bool, int, error)
	
	// GetCount gets the current count for a key
	GetCount(ctx context.Context, key string) (int, error)
	
	// Reset resets the count for a key
	Reset(ctx context.Context, key string) error
}