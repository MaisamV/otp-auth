package entities

import (
	"time"

	"github.com/otp-auth/internal/domain/valueobjects"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID               string                    `json:"id"`
	UserID           string                    `json:"user_id"`
	SessionID        valueobjects.SessionID    `json:"session_id"`
	TokenHash        string                    `json:"token_hash"`
	CreatedAt        time.Time                 `json:"created_at"`
	ExpiresAt        time.Time                 `json:"expires_at"`
	LastUsed         *time.Time                `json:"last_used"`
	Revoked          bool                      `json:"revoked"`
	RevokedAt        *time.Time                `json:"revoked_at"`
	RevokeReason     string                    `json:"revoke_reason"`
}

// RevokeReason constants
const (
	RevokeReasonRefresh = "REFRESH"
	RevokeReasonLogout  = "LOGOUT"
	RevokeReasonExpired = "EXPIRED"
	RevokeReasonAdmin   = "ADMIN"
)

// NewRefreshToken creates a new refresh token
func NewRefreshToken(userID string, sessionID valueobjects.SessionID, tokenHash string, ttl time.Duration) *RefreshToken {
	now := time.Now()
	return &RefreshToken{
		UserID:    userID,
		SessionID: sessionID,
		TokenHash: tokenHash,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
		Revoked:   false,
	}
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the refresh token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.Revoked
}

// Revoke marks the refresh token as revoked with a reason
func (rt *RefreshToken) Revoke(reason string) {
	rt.Revoked = true
	rt.RevokeReason = reason
	now := time.Now()
	rt.RevokedAt = &now
}

// UpdateLastUsed updates the last used timestamp
func (rt *RefreshToken) UpdateLastUsed() {
	now := time.Now()
	rt.LastUsed = &now
}