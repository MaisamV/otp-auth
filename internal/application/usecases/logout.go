package usecases

import (
	"context"

	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/ports/repositories"
	"github.com/otp-auth/internal/application/ports/services"
	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
)

// LogoutUseCase handles user logout operations
type LogoutUseCase struct {
	tokenRepo   repositories.TokenRepository
	hashService services.HashService
}

// NewLogoutUseCase creates a new LogoutUseCase
func NewLogoutUseCase(
	tokenRepo repositories.TokenRepository,
	hashService services.HashService,
) *LogoutUseCase {
	return &LogoutUseCase{
		tokenRepo:   tokenRepo,
		hashService: hashService,
	}
}

// Execute performs the logout operation by revoking the refresh token
func (uc *LogoutUseCase) Execute(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error) {
	// If no refresh token provided, just return success (idempotent operation)
	if req.RefreshToken == "" && req.SessionID == "" {
		return &dto.LogoutResponse{
			Message: "Logout successful",
		}, nil
	}

	// If refresh token is provided, try to revoke it
	if req.RefreshToken != "" {
		// Hash the refresh token to match stored hash
		tokenHash, err := uc.hashService.Hash(req.RefreshToken)
		if err != nil {
			return nil, errors.NewInternalError("Failed to hash refresh token", err)
		}

		// If session ID is also provided, use both for more precise revocation
		if req.SessionID != "" {
			sessionID := valueobjects.SessionID(req.SessionID)
			err = uc.tokenRepo.RevokeByTokenHashAndSessionID(ctx, tokenHash, sessionID, "LOGOUT")
		} else {
			// Revoke by token hash only
			err = uc.tokenRepo.RevokeByTokenHash(ctx, tokenHash, "LOGOUT")
		}

		// If token not found, it's already revoked or expired - still return success
		if err != nil {
			// Check if it's a not found error
			if customErr := errors.GetCustomError(err); customErr != nil && customErr.Type == errors.NotFoundError {
				// Token not found - already revoked or expired, return success
				return &dto.LogoutResponse{
					Message: "Logout successful",
				}, nil
			}
			// Other errors should be returned
			return nil, errors.NewInternalError("Failed to revoke refresh token", err)
		}
	}

	return &dto.LogoutResponse{
		Message: "Logout successful",
	}, nil
}