package dto

import (
	"time"

	"github.com/otp-auth/internal/domain/entities"
)

// SendOTPResponse represents the response after sending an OTP
type SendOTPResponse struct {
	Message   string `json:"message" example:"OTP sent successfully"`
	SessionID string `json:"session_id" example:"abc123def456"`
}

// LoginResponse represents the response after successful login/register
type LoginResponse struct {
	Message          string    `json:"message" example:"Login successful"`
	AccessToken      string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken     string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt        time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at" example:"2024-01-01T12:00:00Z"`
	User             UserInfo  `json:"user"`
}

// RefreshTokenResponse represents the response after token refresh
type RefreshTokenResponse struct {
	Message          string    `json:"message" example:"Token refreshed successfully"`
	AccessToken      string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken     string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt        time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at" example:"2024-01-01T12:00:00Z"`
}

// LogoutResponse represents the response after logout
type LogoutResponse struct {
	Message string `json:"message" example:"Logout successful"`
}

// UserInfo represents user information in responses
type UserInfo struct {
	ID          string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PhoneNumber string    `json:"phone_number" example:"+989123456789"`
	Scope       string    `json:"scope" example:"superadmin"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}

// GetUserResponse represents the response for getting a single user
type GetUserResponse struct {
	User UserInfo `json:"user"`
}

// GetUsersResponse represents the response for getting users list
type GetUsersResponse struct {
	Users      []UserInfo `json:"users"`
	Total      int64      `json:"total" example:"100"`
	Page       int        `json:"page" example:"1"`
	Limit      int        `json:"limit" example:"10"`
	TotalPages int        `json:"total_pages" example:"10"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid phone number format"`
	Code    string `json:"code,omitempty" example:"INVALID_PHONE"`
	Details string `json:"details,omitempty" example:"Phone number must start with + or 0"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// NewUserInfo creates UserInfo from User entity
func NewUserInfo(user *entities.User) UserInfo {
	return UserInfo{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber.String(),
		Scope:       user.Scope,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// NewGetUsersResponse creates GetUsersResponse with pagination
func NewGetUsersResponse(users []*entities.User, total int64, page, limit int) *GetUsersResponse {
	userInfos := make([]UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = NewUserInfo(user)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return &GetUsersResponse{
		Users:      userInfos,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
