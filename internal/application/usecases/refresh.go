package usecases

import (
	"context"
	"log"
	"time"

	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/ports/repositories"
	"github.com/otp-auth/internal/application/ports/services"
	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
)

// RefreshUseCase handles refresh token operations
type RefreshUseCase struct {
	userRepo    repositories.UserRepository
	tokenRepo   repositories.TokenRepository
	jwtService  services.JWTService
	hashService services.HashService
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

// NewRefreshUseCase creates a new RefreshUseCase
func NewRefreshUseCase(
	userRepo repositories.UserRepository,
	tokenRepo repositories.TokenRepository,
	jwtService services.JWTService,
	hashService services.HashService,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *RefreshUseCase {
	return &RefreshUseCase{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		jwtService:  jwtService,
		hashService: hashService,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}

// Execute refreshes the access token using a valid refresh token
func (uc *RefreshUseCase) Execute(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError("Invalid request data", err)
	}

	// Parse session ID
	sessionIDObj, err := valueobjects.NewSessionIDFromString(req.SessionID)
	if err != nil {
		return nil, errors.NewValidationError("Invalid session ID format", err)
	}

	// Hash the provided refresh token to match against stored hash
	hashedRefreshToken, err := uc.hashService.HashRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewInternalError("Failed to hash refresh token", err)
	}

	// Get refresh token from repository
	storedToken, err := uc.tokenRepo.GetByTokenHashAndSessionID(ctx, hashedRefreshToken, sessionIDObj)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid refresh token", err)
	}

	// Validate refresh token
	if !storedToken.IsValid() {
		return nil, errors.NewUnauthorizedError("Refresh token is expired or revoked", nil)
	}

	// Get user information
	user, err := uc.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, errors.NewUnauthorizedError("User not found", err)
	}

	// Convert user scope to scopes array
	var scopes []string
	if user.Scope != "" {
		scopes = []string{user.Scope}
	} else {
		scopes = []string{"user"} // Default scope
	}

	// Generate new access token ID
	accessTokenID, err := uc.hashService.GenerateRandomString(16)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate access token ID", err)
	}

	// Generate new access token claims
	accessClaims := services.NewJWTClaims(
		user.ID,
		"otp-auth-client", // TODO: Make this configurable
		scopes,
		uc.accessTTL,
		"otp-auth", // TODO: Make this configurable
		accessTokenID,
	)

	// Generate new access token
	accessToken, err := uc.jwtService.GenerateToken(accessClaims)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate access token", err)
	}

	// Generate new refresh token
	newRefreshToken, err := uc.hashService.GenerateRandomString(32)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate new refresh token", err)
	}

	// Hash the new refresh token for storage
	hashedNewRefreshToken, err := uc.hashService.HashRefreshToken(newRefreshToken)
	if err != nil {
		return nil, errors.NewInternalError("Failed to hash new refresh token", err)
	}

	// Revoke the old refresh token
	if err := uc.tokenRepo.RevokeByTokenHashAndSessionID(ctx, hashedRefreshToken, sessionIDObj, entities.RevokeReasonRefresh); err != nil {
		// Log error but don't fail the refresh operation
		log.Printf("revoke error: %v", err)
		// TODO: Add proper logging
	}

	// Create new refresh token entity
	newRefreshTokenEntity := entities.NewRefreshToken(
		user.ID,
		sessionIDObj,
		hashedNewRefreshToken,
		uc.refreshTTL,
	)
	newRefreshTokenEntity.ID = generateTokenID() // Generate unique ID

	// Store new refresh token
	if err := uc.tokenRepo.Create(ctx, newRefreshTokenEntity); err != nil {
		return nil, errors.NewInternalError("Failed to store new refresh token", err)
	}

	// Update last used timestamp for the old token (before revocation)
	storedToken.UpdateLastUsed()
	if err := uc.tokenRepo.Update(ctx, storedToken); err != nil {
		// Log error but don't fail the refresh operation
		// TODO: Add proper logging
	}

	// Calculate access token expiration
	expiresAt := time.Now().Add(uc.accessTTL)
	refreshExpiresAt := time.Now().Add(uc.refreshTTL)

	// Create response
	response := &dto.RefreshTokenResponse{
		Message:          "Token refreshed successfully",
		AccessToken:      accessToken,
		RefreshToken:     newRefreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}

	return response, nil
}
