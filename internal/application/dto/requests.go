package dto

import (
	"github.com/otp-auth/internal/domain/valueobjects"
)

// SendOTPRequest represents the request to send an OTP
type SendOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"+989123456789"`
	SessionID   string `json:"session_id,omitempty" example:"abc123def456"`
}

// LoginRequest represents the request to login/register
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"+989123456789"`
	OTP         string `json:"otp" binding:"required" example:"123456"`
}

// RefreshTokenRequest represents the request to refresh tokens
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	SessionID    string `json:"-"` // This field is populated from cookies, not request body
}

// LogoutRequest represents the request to logout
// Both RefreshToken and SessionID are populated from cookies, not request body
type LogoutRequest struct {
	RefreshToken string `json:"-"` // Read from cookies
	SessionID    string `json:"-"` // Read from cookies
}

// GetUsersRequest represents the request to get users list
type GetUsersRequest struct {
	Page         int    `form:"page" example:"1"`
	Limit        int    `form:"limit" example:"10"`
	SearchPhone  string `form:"search_phone" example:"+989123456789"`
	DateFrom     string `form:"date_from" example:"2024-01-01"`
	DateTo       string `form:"date_to" example:"2024-12-31"`
}

// UpdateUserScopeRequest represents the request to update user scope
type UpdateUserScopeRequest struct {
	Scope string `json:"scope" binding:"required" example:"superadmin"`
}

// Validate validates the SendOTPRequest
func (r *SendOTPRequest) Validate() error {
	_, err := valueobjects.NewPhoneNumber(r.PhoneNumber)
	return err
}

// Validate validates the LoginRequest
func (r *LoginRequest) Validate() error {
	if _, err := valueobjects.NewPhoneNumber(r.PhoneNumber); err != nil {
		return err
	}
	
	return nil
}

// Validate validates the RefreshTokenRequest
func (r *RefreshTokenRequest) Validate() error {
	// SessionID validation is handled separately since it comes from cookies
	if r.SessionID != "" {
		_, err := valueobjects.NewSessionIDFromString(r.SessionID)
		return err
	}
	return nil
}