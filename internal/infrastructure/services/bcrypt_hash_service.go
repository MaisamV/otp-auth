package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/otp-auth/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

// BcryptHashService implements HashService using bcrypt algorithm
type BcryptHashService struct {
	cost int
}

// HashConfig holds configuration for hash service
type HashConfig struct {
	Cost int // bcrypt cost factor (4-31, recommended: 10-12)
}

// DefaultHashConfig returns default hash configuration
func DefaultHashConfig() HashConfig {
	return HashConfig{
		Cost: bcrypt.DefaultCost, // Usually 10
	}
}

// NewBcryptHashService creates a new bcrypt hash service
func NewBcryptHashService(config HashConfig) *BcryptHashService {
	cost := config.Cost
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	return &BcryptHashService{
		cost: cost,
	}
}

// HashPassword hashes a password using bcrypt
func (h *BcryptHashService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.NewValidationError("Password cannot be empty", nil)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", errors.NewInternalError("Failed to hash password", err)
	}

	return string(hash), nil
}

// VerifyPassword verifies a password against its hash
func (h *BcryptHashService) VerifyPassword(password, hash string) error {
	if password == "" {
		return errors.NewValidationError("Password cannot be empty", nil)
	}

	if hash == "" {
		return errors.NewValidationError("Hash cannot be empty", nil)
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("error: %v", err)
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.NewValidationError("Invalid password", nil)
		}
		return errors.NewInternalError("Failed to verify password", err)
	}

	return nil
}

// HashOTP hashes an OTP code using bcrypt
func (h *BcryptHashService) HashOTP(otp string) (string, error) {
	if otp == "" {
		return "", errors.NewValidationError("OTP cannot be empty", nil)
	}

	// Use a lower cost for OTP hashing since they're short-lived
	// and we need faster verification
	otpCost := h.cost
	if otpCost > 8 {
		otpCost = 8 // Max cost of 8 for OTPs
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(otp), otpCost)
	if err != nil {
		return "", errors.NewInternalError("Failed to hash OTP", err)
	}

	return string(hash), nil
}

// VerifyOTP verifies an OTP code against its hash
func (h *BcryptHashService) VerifyOTP(otp, hash string) error {
	if otp == "" {
		return errors.NewValidationError("OTP cannot be empty", nil)
	}

	if hash == "" {
		return errors.NewValidationError("Hash cannot be empty", nil)
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(otp))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.NewValidationError("Invalid OTP", nil)
		}
		return errors.NewInternalError("Failed to verify OTP", err)
	}

	return nil
}

// HashRefreshToken hashes a refresh token using SHA-256
func (h *BcryptHashService) HashRefreshToken(token string) (string, error) {
	if token == "" {
		return "", errors.NewValidationError("Token cannot be empty", nil)
	}

	hasher := sha256.New()
	hasher.Write([]byte(token))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash, nil
}

// VerifyRefreshToken verifies a refresh token against its hash
func (h *BcryptHashService) VerifyRefreshToken(token, hash string) error {
	if token == "" {
		return errors.NewValidationError("Token cannot be empty", nil)
	}

	if hash == "" {
		return errors.NewValidationError("Hash cannot be empty", nil)
	}

	hasher := sha256.New()
	hasher.Write([]byte(token))
	computedHash := hex.EncodeToString(hasher.Sum(nil))

	if computedHash != hash {
		return errors.NewValidationError("Invalid refresh token", nil)
	}

	return nil
}

// GetCost returns the current bcrypt cost
func (h *BcryptHashService) GetCost() int {
	return h.cost
}

// SetCost updates the bcrypt cost (for runtime configuration)
func (h *BcryptHashService) SetCost(cost int) error {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return errors.NewValidationError("Invalid bcrypt cost: must be between 4 and 31", nil)
	}

	h.cost = cost
	return nil
}

// GetHashInfo returns information about the hash service
func (h *BcryptHashService) GetHashInfo() map[string]interface{} {
	return map[string]interface{}{
		"algorithm": "bcrypt",
		"cost":      h.cost,
	}
}

// Hash implements the HashService interface
func (h *BcryptHashService) Hash(plaintext string) (string, error) {
	return h.HashPassword(plaintext)
}

// Compare implements the HashService interface
func (h *BcryptHashService) Compare(plaintext, hash string) bool {
	return h.VerifyPassword(plaintext, hash) == nil
}

// GenerateRandomString generates a random string of specified length
func (h *BcryptHashService) GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}
