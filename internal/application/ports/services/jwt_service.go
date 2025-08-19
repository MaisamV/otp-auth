package services

import (
	"time"
)

// JWTClaims represents the structure of JWT claims
type JWTClaims struct {
	Subject   string   `json:"sub"`       // User ID
	ClientID  string   `json:"client_id"` // Client identifier
	Scopes    []string `json:"scopes"`    // User scopes/permissions
	IssuedAt  int64    `json:"iat"`       // Issued at timestamp
	ExpiresAt int64    `json:"exp"`       // Expiration timestamp
	Issuer    string   `json:"iss"`       // Token issuer
	TokenID   string   `json:"jti"`       // JWT ID (unique token identifier)
}

// NewJWTClaims creates new JWT claims with the given parameters
func NewJWTClaims(userID, clientID string, scopes []string, ttl time.Duration, issuer, tokenID string) *JWTClaims {
	now := time.Now()
	return &JWTClaims{
		Subject:   userID,
		ClientID:  clientID,
		Scopes:    scopes,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(ttl).Unix(),
		Issuer:    issuer,
		TokenID:   tokenID,
	}
}

// IsExpired checks if the token has expired
func (c *JWTClaims) IsExpired() bool {
	return time.Now().Unix() > c.ExpiresAt
}

// JWTService defines the interface for JWT token operations
type JWTService interface {
	// GenerateToken generates a JWT token from claims
	GenerateToken(claims *JWTClaims) (string, error)
	
	// VerifyToken verifies and parses a JWT token, returning the claims
	VerifyToken(token string) (*JWTClaims, error)
}