package repositories

import (
	"context"

	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
)

// TokenReader defines read operations for refresh tokens
type TokenReader interface {
	// GetByTokenHash retrieves a refresh token by token hash
	GetByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error)
	
	// GetByTokenHashAndSessionID retrieves a refresh token by token hash and session ID
	GetByTokenHashAndSessionID(ctx context.Context, tokenHash string, sessionID valueobjects.SessionID) (*entities.RefreshToken, error)
	
	// GetByUserID retrieves all refresh tokens for a user
	GetByUserID(ctx context.Context, userID string) ([]*entities.RefreshToken, error)
	
	// GetActiveByUserID retrieves all active (non-revoked, non-expired) refresh tokens for a user
	GetActiveByUserID(ctx context.Context, userID string) ([]*entities.RefreshToken, error)
}

// TokenWriter defines write operations for refresh tokens
type TokenWriter interface {
	// Create creates a new refresh token
	Create(ctx context.Context, token *entities.RefreshToken) error
	
	// Update updates an existing refresh token
	Update(ctx context.Context, token *entities.RefreshToken) error
	
	// RevokeByTokenHash revokes a refresh token by token hash
	RevokeByTokenHash(ctx context.Context, tokenHash string, reason string) error
	
	// RevokeByTokenHashAndSessionID revokes a refresh token by token hash and session ID
	RevokeByTokenHashAndSessionID(ctx context.Context, tokenHash string, sessionID valueobjects.SessionID, reason string) error
	
	// RevokeAllByUserID revokes all refresh tokens for a user
	RevokeAllByUserID(ctx context.Context, userID string, reason string) error
	
	// Delete deletes a refresh token by ID
	Delete(ctx context.Context, id string) error
	
	// DeleteExpired deletes all expired refresh tokens
	DeleteExpired(ctx context.Context) error
}

// TokenRepository combines read and write operations
type TokenRepository interface {
	TokenReader
	TokenWriter
}