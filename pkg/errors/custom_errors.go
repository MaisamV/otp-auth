package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ValidationError ErrorType = "VALIDATION_ERROR"
	NotFoundError   ErrorType = "NOT_FOUND_ERROR"
	Unauthorized    ErrorType = "UNAUTHORIZED_ERROR"
	Forbidden       ErrorType = "FORBIDDEN_ERROR"
	RateLimitError  ErrorType = "RATE_LIMIT_ERROR"
	InternalError   ErrorType = "INTERNAL_ERROR"
	ConflictError   ErrorType = "CONFLICT_ERROR"
	GoneError       ErrorType = "GONE_ERROR"
)

// CustomError represents a custom application error
type CustomError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"status_code"`
	Cause      error     `json:"-"`
}

// Error implements the error interface
func (e *CustomError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *CustomError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a new validation error
func NewValidationError(message string, cause error) *CustomError {
	details := ""
	if cause != nil {
		details = cause.Error()
	}
	return &CustomError{
		Type:       ValidationError,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
		Cause:      cause,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, cause error) *CustomError {
	details := ""
	if cause != nil {
		details = cause.Error()
	}
	return &CustomError{
		Type:       NotFoundError,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusNotFound,
		Cause:      cause,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, cause error) *CustomError {
	details := ""
	if cause != nil {
		details = cause.Error()
	}
	return &CustomError{
		Type:       Unauthorized,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusUnauthorized,
		Cause:      cause,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, cause error) *CustomError {
	details := ""
	if cause != nil {
		details = cause.Error()
	}
	return &CustomError{
		Type:       Forbidden,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusForbidden,
		Cause:      cause,
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string) *CustomError {
	return &CustomError{
		Type:       RateLimitError,
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string, cause error) *CustomError {
	details := ""
	if cause != nil {
		details = cause.Error()
	}
	return &CustomError{
		Type:       InternalError,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusInternalServerError,
		Cause:      cause,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string, cause error) *CustomError {
	details := ""
	if cause != nil {
		details = cause.Error()
	}
	return &CustomError{
		Type:       ConflictError,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusConflict,
		Cause:      cause,
	}
}

// NewGoneError creates a new gone error (for expired resources)
func NewGoneError(message string, cause error) *CustomError {
	details := ""
	if cause != nil {
		details = cause.Error()
	}
	return &CustomError{
		Type:       GoneError,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusGone,
		Cause:      cause,
	}
}

// IsCustomError checks if an error is a CustomError
func IsCustomError(err error) bool {
	_, ok := err.(*CustomError)
	return ok
}

// GetCustomError extracts CustomError from error
func GetCustomError(err error) *CustomError {
	if customErr, ok := err.(*CustomError); ok {
		return customErr
	}
	return nil
}