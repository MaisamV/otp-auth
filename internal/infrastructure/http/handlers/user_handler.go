package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/usecases"
	"github.com/otp-auth/internal/infrastructure/http/middleware"
	"github.com/otp-auth/pkg/errors"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	getUserProfileUseCase *usecases.GetUserProfileUseCase
	getUsersListUseCase   *usecases.GetUsersListUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(getUserProfileUseCase *usecases.GetUserProfileUseCase, getUsersListUseCase *usecases.GetUsersListUseCase) *UserHandler {
	return &UserHandler{
		getUserProfileUseCase: getUserProfileUseCase,
		getUsersListUseCase:   getUsersListUseCase,
	}
}

// GetUsers handles the get users request (admin only)
// @Summary Get Users
// @Description Get list of users with pagination (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search_phone query string false "Search by phone number"
// @Param date_from query string false "Filter from date (YYYY-MM-DD)"
// @Param date_to query string false "Filter to date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} dto.GetUsersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	searchPhone := c.Query("search_phone")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	req := &dto.GetUsersRequest{
		Page:        page,
		Limit:       limit,
		SearchPhone: searchPhone,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	// Execute get users list use case
	response, err := h.getUsersListUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserScope handles the update user scope request (admin only)
// @Summary Update User Scope
// @Description Update user scope (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserScopeRequest true "Update user scope request"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id}/scope [put]
func (h *UserHandler) UpdateUserScope(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		h.handleError(c, errors.NewValidationError("User ID is required", nil))
		return
	}

	var req dto.UpdateUserScopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.NewValidationError("Invalid request format", err))
		return
	}

	// TODO: Implement update user scope use case
	// For now, return placeholder response
	_ = userID // Use the variable to avoid "declared and not used" error
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Update user scope endpoint not implemented yet"})
}

// GetProfile handles the get user profile request
// @Summary Get User Profile
// @Description Get current user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.GetUserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, exists := middleware.GetUserID(c)
	if !exists {
		h.handleError(c, errors.NewUnauthorizedError("User ID not found in context", nil))
		return
	}

	// Execute get user profile use case
	response, err := h.getUserProfileUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// handleError handles errors and sends appropriate HTTP responses
func (h *UserHandler) handleError(c *gin.Context, err error) {
	if customErr, ok := err.(*errors.CustomError); ok {
		c.JSON(customErr.StatusCode, dto.ErrorResponse{
			Error:   customErr.Message,
			Code:    string(customErr.Type),
			Details: customErr.Details,
		})
		return
	}

	// Default to internal server error
	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Error:   "An internal error occurred",
		Code:    "INTERNAL_ERROR",
		Details: err.Error(),
	})
}