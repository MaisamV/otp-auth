package usecases

import (
	"context"

	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/ports/repositories"
)

// GetUserProfileUseCase handles getting user profile operations
type GetUserProfileUseCase struct {
	userRepo repositories.UserRepository
}

// NewGetUserProfileUseCase creates a new GetUserProfileUseCase
func NewGetUserProfileUseCase(userRepo repositories.UserRepository) *GetUserProfileUseCase {
	return &GetUserProfileUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves the user profile by user ID
func (uc *GetUserProfileUseCase) Execute(ctx context.Context, userID string) (*dto.GetUserResponse, error) {
	// Get user by ID
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	userInfo := dto.NewUserInfo(user)
	response := &dto.GetUserResponse{
		User: userInfo,
	}

	return response, nil
}