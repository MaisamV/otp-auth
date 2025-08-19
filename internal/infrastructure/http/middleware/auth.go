package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/otp-auth/internal/application/dto"
	"github.com/otp-auth/internal/application/ports/services"
	"net/http"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtService services.JWTService
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(jwtService services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth middleware that requires valid JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			m.unauthorizedResponse(c, "Authorization token required")
			return
		}

		claims, err := m.jwtService.VerifyToken(token)
		if err != nil {
			m.unauthorizedResponse(c, "Invalid or expired token")
			return
		}

		// Check if token is expired
		if claims.IsExpired() {
			m.unauthorizedResponse(c, "Token has expired")
			return
		}

		// Set user information in context
		c.Set("user_id", claims.Subject)
		c.Set("client_id", claims.ClientID)
		c.Set("scopes", claims.Scopes)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireAdmin middleware that requires admin scope
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			m.unauthorizedResponse(c, "Authorization token required")
			return
		}

		claims, err := m.jwtService.VerifyToken(token)
		if err != nil {
			m.unauthorizedResponse(c, "Invalid or expired token")
			return
		}

		// Check if token is expired
		if claims.IsExpired() {
			m.unauthorizedResponse(c, "Token has expired")
			return
		}

		// Check if user has admin scope
		hasAdminScope := false
		for _, scope := range claims.Scopes {
			if scope == "admin" || scope == "superadmin" {
				hasAdminScope = true
				break
			}
		}

		if !hasAdminScope {
			m.forbiddenResponse(c, "Admin access required")
			return
		}

		// Set user information in context
		c.Set("user_id", claims.Subject)
		c.Set("client_id", claims.ClientID)
		c.Set("scopes", claims.Scopes)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuth middleware that optionally validates JWT token
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := m.jwtService.VerifyToken(token)
		if err != nil {
			// Invalid token, but continue without authentication
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.Subject)
		c.Set("client_id", claims.ClientID)
		c.Set("scopes", claims.Scopes)
		c.Set("claims", claims)

		c.Next()
	}
}

// extractToken extracts JWT token from access_token cookie
func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	accessToken, _ := c.Cookie("access_token")
	return accessToken
}

// unauthorizedResponse sends an unauthorized response
func (m *AuthMiddleware) unauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
		Error: message,
		Code:  "UNAUTHORIZED",
	})
	c.Abort()
}

// forbiddenResponse sends a forbidden response
func (m *AuthMiddleware) forbiddenResponse(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, dto.ErrorResponse{
		Error: message,
		Code:  "FORBIDDEN",
	})
	c.Abort()
}

// Helper functions to get user information from context

// GetUserID gets the user ID from the request context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// GetClientID gets the client ID from the request context
func GetClientID(c *gin.Context) (string, bool) {
	clientID, exists := c.Get("client_id")
	if !exists {
		return "", false
	}

	clientIDStr, ok := clientID.(string)
	return clientIDStr, ok
}

// GetScopes gets the scopes from the request context
func GetScopes(c *gin.Context) ([]string, bool) {
	scopes, exists := c.Get("scopes")
	if !exists {
		return nil, false
	}

	scopesSlice, ok := scopes.([]string)
	return scopesSlice, ok
}

// GetClaims gets the JWT claims from the request context
func GetClaims(c *gin.Context) (*services.JWTClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*services.JWTClaims)
	return jwtClaims, ok
}
