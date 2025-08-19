package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/ports/repositories"
	"github.com/otp-auth/internal/application/ports/services"
	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
)

// LoginUseCase handles user login/registration with OTP verification
type LoginUseCase struct {
	userRepo    repositories.UserRepository
	otpRepo     repositories.OTPRepository
	tokenRepo   repositories.TokenRepository
	jwtService  services.JWTService
	hashService services.HashService
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

// NewLoginUseCase creates a new LoginUseCase
func NewLoginUseCase(
	userRepo repositories.UserRepository,
	otpRepo repositories.OTPRepository,
	tokenRepo repositories.TokenRepository,
	jwtService services.JWTService,
	hashService services.HashService,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		otpRepo:     otpRepo,
		tokenRepo:   tokenRepo,
		jwtService:  jwtService,
		hashService: hashService,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}

// Execute performs the login/registration process
func (uc *LoginUseCase) Execute(ctx context.Context, req *dto.LoginRequest, sessionID string) (*dto.LoginResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError("Invalid request", err)
	}

	// Convert phone number string to PhoneNumber value object
	phoneNumber, err := valueobjects.NewPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, errors.NewValidationError("Invalid phone number format", err)
	}

	// Get OTP from repository
	storedOTP, err := uc.otpRepo.Get(ctx, phoneNumber)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid session ID", err)
	}

	// Verify OTP
	if err := uc.hashService.VerifyOTP(req.OTP, storedOTP.HashedCode); err != nil {
		return nil, errors.NewUnauthorizedError("Invalid OTP", err)
	}

	// Note: OTP expiration is handled by Redis TTL, so if we can retrieve it, it's still valid

	// Convert session ID string to SessionID value object
	sessionIDObj, err := valueobjects.NewSessionIDFromString(sessionID)
	if err != nil {
		return nil, errors.NewValidationError("Invalid session ID format", err)
	}

	// Validate session ID matches
	if storedOTP.SessionID != sessionIDObj {
		return nil, errors.NewUnauthorizedError("Session ID mismatch", nil)
	}

	// Delete used OTP
	if err := uc.otpRepo.Delete(ctx, phoneNumber); err != nil {
		// Log error but don't fail the login
		// TODO: Add proper logging
	}

	// Check if user exists
	user, err := uc.userRepo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		// User doesn't exist, create new user (registration)
		user = entities.NewUser(phoneNumber)
		user.ID = generateUserID() // Generate unique ID

		if err := uc.userRepo.Create(ctx, user); err != nil {
			return nil, errors.NewInternalError("Failed to create user", err)
		}
	}

	// Convert user scope to scopes array
	var scopes []string
	if user.Scope != "" {
		scopes = []string{user.Scope}
	} else {
		scopes = []string{"user"} // Default scope
	}

	// Generate token ID for access token
	accessTokenID, err := uc.hashService.GenerateRandomString(16)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate access token ID", err)
	}

	// Generate access token claims
	accessClaims := services.NewJWTClaims(
		user.ID,
		"otp-auth-client", // TODO: Make this configurable
		scopes,
		uc.accessTTL, // Access token TTL from config
		"otp-auth",   // TODO: Make this configurable
		accessTokenID,
	)

	// Generate access token
	fmt.Println("[DEBUG] About to call GenerateToken in login use case")
	accessToken, err := uc.jwtService.GenerateToken(accessClaims)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate access token", err)
	}

	// Generate refresh token as a random string (not JWT)
	refreshToken, err := uc.hashService.GenerateRandomString(32)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate refresh token", err)
	}

	// Hash the refresh token for storage
	hashedRefreshToken, err := uc.hashService.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.NewInternalError("Failed to hash refresh token", err)
	}

	// Store refresh token in repository
	refreshTokenEntity := entities.NewRefreshToken(
		user.ID,
		sessionIDObj,
		hashedRefreshToken,
		uc.refreshTTL,
	)
	refreshTokenEntity.ID = generateTokenID() // Generate unique ID

	if err := uc.tokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		return nil, errors.NewInternalError("Failed to store refresh token", err)
	}

	// Calculate access token expiration
	now := time.Now()
	expiresAt := now.Add(uc.accessTTL)
	refreshExpiresAt := now.Add(uc.refreshTTL)

	// Create response
	response := &dto.LoginResponse{
		Message:          "Login successful",
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		User:             dto.NewUserInfo(user),
	}

	return response, nil
}

// generateUserID generates a unique user ID
func generateUserID() string {
	return uuid.New().String()
}

// generateTokenID generates a unique token ID
func generateTokenID() string {
	return uuid.New().String()
}
