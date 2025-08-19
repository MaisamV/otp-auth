package usecases

import (
	"context"
	"time"

	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/ports/repositories"
	"github.com/otp-auth/pkg/errors"
)

// GetUsersListUseCase handles getting users list operations
type GetUsersListUseCase struct {
	userRepo repositories.UserRepository
}

// NewGetUsersListUseCase creates a new GetUsersListUseCase
func NewGetUsersListUseCase(userRepo repositories.UserRepository) *GetUsersListUseCase {
	return &GetUsersListUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves users list with pagination and search
func (uc *GetUsersListUseCase) Execute(ctx context.Context, req *dto.GetUsersRequest) (*dto.GetUsersResponse, error) {
	// Validate pagination parameters
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // Maximum limit
	}

	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Validate date formats if provided
	var dateFrom, dateTo string
	if req.DateFrom != "" {
		if _, err := time.Parse("2006-01-02", req.DateFrom); err != nil {
			return nil, errors.NewValidationError("Invalid date_from format. Use YYYY-MM-DD", err)
		}
		dateFrom = req.DateFrom + " 00:00:00"
	}
	if req.DateTo != "" {
		if _, err := time.Parse("2006-01-02", req.DateTo); err != nil {
			return nil, errors.NewValidationError("Invalid date_to format. Use YYYY-MM-DD", err)
		}
		dateTo = req.DateTo + " 23:59:59"
	}

	// Get users from repository
	users, total, err := uc.userRepo.List(ctx, offset, req.Limit, req.SearchPhone, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Create response with pagination
	response := dto.NewGetUsersResponse(users, total, req.Page, req.Limit)
	return response, nil
}