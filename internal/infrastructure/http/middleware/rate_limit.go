package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/otp-auth/internal/application/ports/repositories"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/otp-auth/internal/application/dto"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Limit    int                       // Number of requests allowed
	Window   time.Duration             // Time window
	KeyFunc  func(*gin.Context) string // Function to generate rate limit key
	SkipFunc func(*gin.Context) bool   // Function to skip rate limiting
}

// DefaultRateLimitConfig returns a default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Limit:  100,
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		SkipFunc: func(c *gin.Context) bool {
			return false
		},
	}
}

// RateLimit returns a rate limiting middleware
func RateLimit(rateLimiter repositories.RateLimiter, config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting if skip function returns true
		if config.SkipFunc(c) {
			c.Next()
			return
		}

		// Generate rate limit key
		key := config.KeyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		// Check rate limit and increment counter
		allowed, count, err := rateLimiter.CheckAndIncrement(c.Request.Context(), key, config.Limit, config.Window)
		if err != nil {
			// Log error but don't block request
			c.Next()
			return
		}

		if !allowed {
			// Set rate limit headers
			c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Window).Unix(), 10))
			c.Header("X-RateLimit-Window", config.Window.String())

			// Return rate limit error
			detailsMap := map[string]interface{}{
				"limit":    config.Limit,
				"window":   config.Window.String(),
				"current":  count,
				"reset_at": time.Now().Add(config.Window).Unix(),
			}
			detailsJSON, _ := json.Marshal(detailsMap)
			errorResp := dto.ErrorResponse{
				Error:   "Rate limit exceeded",
				Code:    "RATE_LIMIT_ERROR",
				Details: string(detailsJSON),
			}
			c.JSON(http.StatusTooManyRequests, errorResp)
			c.Abort()
			return
		}
		remaining := config.Limit - count
		if remaining < 0 {
			remaining = 0
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Window).Unix(), 10))
		c.Header("X-RateLimit-Window", config.Window.String())

		c.Next()
	}
}

// IPBasedRateLimit creates a rate limit middleware based on client IP
func IPBasedRateLimit(rateLimiter repositories.RateLimiter, limit int, window time.Duration) gin.HandlerFunc {
	config := RateLimitConfig{
		Limit:  limit,
		Window: window,
		KeyFunc: func(c *gin.Context) string {
			return fmt.Sprintf("ip:%s", c.ClientIP())
		},
		SkipFunc: func(c *gin.Context) bool {
			return false
		},
	}
	return RateLimit(rateLimiter, config)
}

// OtpRateLimit creates a rate limit middleware based on phone number
func OtpRateLimit(rateLimiter repositories.RateLimiter, limit int, window time.Duration) gin.HandlerFunc {
	config := RateLimitConfig{
		Limit:  limit,
		Window: window,
		KeyFunc: func(c *gin.Context) string {
			// Try to get phone number from request body
			var req struct {
				PhoneNumber string `json:"phone_number"`
			}

			// Read the request body
			body, err := c.GetRawData()
			if err != nil {
				return ""
			}

			// Restore the request body for the next handler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// Parse JSON to get phone number
			if err := json.Unmarshal(body, &req); err == nil && req.PhoneNumber != "" {
				return fmt.Sprintf("phone:%s", req.PhoneNumber)
			}

			return ""
		},
		SkipFunc: func(c *gin.Context) bool {
			return false
		},
	}
	return RateLimit(rateLimiter, config)
}
