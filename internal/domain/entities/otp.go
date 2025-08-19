package entities

import (
	"time"

	"github.com/otp-auth/internal/domain/valueobjects"
)

// OTP represents an OTP code in the system
type OTP struct {
	PhoneNumber valueobjects.PhoneNumber `json:"phone_number"`
	SessionID   valueobjects.SessionID   `json:"session_id"`
	HashedCode  string                   `json:"hashed_code"`
	CreatedAt   time.Time                `json:"created_at"`
	ExpiresAt   time.Time                `json:"expires_at"`
}

// NewOTP creates a new OTP with the given parameters
func NewOTP(phoneNumber valueobjects.PhoneNumber, sessionID valueobjects.SessionID, hashedCode string, ttl time.Duration) *OTP {
	now := time.Now()
	return &OTP{
		PhoneNumber: phoneNumber,
		SessionID:   sessionID,
		HashedCode:  hashedCode,
		CreatedAt:   now,
		ExpiresAt:   now.Add(ttl),
	}
}

// IsExpired checks if the OTP has expired
func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsValid checks if the OTP is valid (not expired)
func (o *OTP) IsValid() bool {
	return !o.IsExpired()
}
