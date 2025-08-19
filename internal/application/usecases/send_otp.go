package usecases

import (
	"context"
	"time"

	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/ports/repositories"
	"github.com/otp-auth/internal/application/ports/services"
	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
	"github.com/otp-auth/pkg/utils"
)

// SendOTPUseCase handles the send OTP flow
type SendOTPUseCase struct {
	userRepo        repositories.UserRepository
	otpRepo         repositories.OTPRepository
	rateLimiter     repositories.RateLimiter
	otpSender       services.OTPSender
	hashService     services.HashService
	otpTTL          time.Duration
	rateLimitWindow time.Duration
	rateLimitMax    int
}

// NewSendOTPUseCase creates a new SendOTPUseCase
func NewSendOTPUseCase(userRepo repositories.UserRepository, otpRepo repositories.OTPRepository, rateLimiter repositories.RateLimiter, otpSender services.OTPSender, hashService services.HashService, otpTTL time.Duration, rateLimitWindow time.Duration, rateLimitMax int) *SendOTPUseCase {
	return &SendOTPUseCase{
		userRepo:        userRepo,
		otpRepo:         otpRepo,
		rateLimiter:     rateLimiter,
		otpSender:       otpSender,
		hashService:     hashService,
		otpTTL:          otpTTL,
		rateLimitWindow: rateLimitWindow,
		rateLimitMax:    rateLimitMax,
	}
}

// Execute executes the send OTP use case
func (uc *SendOTPUseCase) Execute(ctx context.Context, req *dto.SendOTPRequest) (*dto.SendOTPResponse, error) {
	// Validate phone number
	phoneNumber, err := valueobjects.NewPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, errors.NewValidationError("Invalid phone number format", err)
	}

	// Get or create session ID
	var sessionID valueobjects.SessionID
	if req.SessionID != "" {
		sessionID, err = valueobjects.NewSessionIDFromString(req.SessionID)
		if err != nil {
			return nil, errors.NewValidationError("Invalid session ID format", err)
		}
	} else {
		sessionID, err = valueobjects.NewSessionID()
		if err != nil {
			return nil, errors.NewInternalError("Failed to generate session ID", err)
		}
	}

	// Generate OTP
	otpCode, err := uc.generateOTP()
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate OTP", err)
	}

	// Hash OTP for security
	hashedOTP, err := uc.hashService.HashOTP(otpCode)
	if err != nil {
		return nil, errors.NewInternalError("Failed to hash OTP", err)
	}
	// Create OTP entity
	otpEntity := entities.NewOTP(phoneNumber, sessionID, hashedOTP, uc.otpTTL)

	// Store OTP in Redis
	if err := uc.otpRepo.Store(ctx, otpEntity, uc.otpTTL); err != nil {
		return nil, errors.NewInternalError("Failed to store OTP", err)
	}

	// Send OTP
	if err := uc.otpSender.SendOTP(ctx, phoneNumber, otpCode); err != nil {
		return nil, errors.NewInternalError("Failed to send OTP", err)
	}

	return &dto.SendOTPResponse{
		Message:   "OTP sent successfully",
		SessionID: sessionID.String(),
	}, nil
}

// generateOTP generates a 6-digit OTP
func (uc *SendOTPUseCase) generateOTP() (string, error) {
	// Use the proper random OTP generation from utils
	return utils.GenerateOTP(6)
}
