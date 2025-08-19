package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/usecases"
	"github.com/otp-auth/pkg/errors"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	sendOTPUseCase *usecases.SendOTPUseCase
	loginUseCase   *usecases.LoginUseCase
	refreshUseCase *usecases.RefreshUseCase
	logoutUseCase  *usecases.LogoutUseCase
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(sendOTPUseCase *usecases.SendOTPUseCase, loginUseCase *usecases.LoginUseCase, refreshUseCase *usecases.RefreshUseCase, logoutUseCase *usecases.LogoutUseCase) *AuthHandler {
	return &AuthHandler{
		sendOTPUseCase: sendOTPUseCase,
		loginUseCase:   loginUseCase,
		refreshUseCase: refreshUseCase,
		logoutUseCase:  logoutUseCase,
	}
}

// SendOTP handles the send OTP request
// @Summary Send OTP
// @Description Send OTP to phone number for authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.SendOTPRequest true "Send OTP request"
// @Success 200 {object} dto.SendOTPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 429 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/send-otp [post]
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.NewValidationError("Invalid request format", err))
		return
	}

	// Check for existing session_id in cookies
	if sessionIDCookie, err := c.Cookie("session_id"); err == nil && sessionIDCookie != "" {
		// Use existing session_id from cookie
		req.SessionID = sessionIDCookie
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.handleError(c, errors.NewValidationError("Invalid phone number format", err))
		return
	}

	// Execute use case
	response, err := h.sendOTPUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Set session ID as httpOnly cookie
	c.SetCookie("session_id", response.SessionID, 259200, "/", "", false, true)

	c.JSON(http.StatusOK, response)
}

// Login handles the login/register request
// @Summary Login/Register
// @Description Login or register user with OTP verification
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.NewValidationError("Invalid request format", err))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.handleError(c, errors.NewValidationError("Invalid request data", err))
		return
	}

	// Get session ID from cookie
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		h.handleError(c, errors.NewUnauthorizedError("Session ID not found in cookies", err))
		return
	}

	// Execute login use case
	response, err := h.loginUseCase.Execute(c.Request.Context(), &req, sessionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Set access token as HTTP-only cookie
	c.SetCookie("access_token", response.AccessToken, int(response.ExpiresAt.Sub(time.Now()).Seconds()), "/", "", false, true)

	// Set refresh token as HTTP-only cookie
	c.SetCookie("refresh_token", response.RefreshToken, int(response.RefreshExpiresAt.Sub(time.Now()).Seconds()), "/", "", false, true)

	// Set session ID as HTTP-only cookie
	c.SetCookie("session_id", sessionID, int(response.RefreshExpiresAt.Sub(time.Now()).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles the refresh token request
// @Summary Refresh Token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.handleError(c, errors.NewUnauthorizedError("Refresh token not found in cookies", err))
		return
	}

	sessionID, err := c.Cookie("session_id")
	if err != nil {
		h.handleError(c, errors.NewUnauthorizedError("Session ID not found in cookies", err))
		return
	}

	req := dto.RefreshTokenRequest{
		refreshToken,
		sessionID,
	}
	// Validate request
	if err := req.Validate(); err != nil {
		h.handleError(c, errors.NewValidationError("Invalid request data", err))
		return
	}

	// Execute refresh use case
	response, err := h.refreshUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Set new access token as HTTP-only cookie
	c.SetCookie("access_token", response.AccessToken, int(response.ExpiresAt.Sub(time.Now()).Seconds()), "/", "", false, true)

	// Set new refresh token as HTTP-only cookie
	c.SetCookie("refresh_token", response.RefreshToken, int(response.RefreshExpiresAt.Sub(time.Now()).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, response)
}

// Logout handles the logout request
// @Summary Logout
// @Description Logout user and invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "Logout request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Read refresh token and session ID from cookies
	refreshToken, _ := c.Cookie("refresh_token")
	sessionID, _ := c.Cookie("session_id")

	// Create logout request with data from cookies
	req := dto.LogoutRequest{
		RefreshToken: refreshToken,
		SessionID:    sessionID,
	}

	// Execute logout use case
	response, err := h.logoutUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Clear all authentication cookies
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, response)
}

// handleError handles errors and sends appropriate HTTP responses
func (h *AuthHandler) handleError(c *gin.Context, err error) {
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
