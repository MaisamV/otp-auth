package router

import (
	"github.com/gin-gonic/gin"
	"github.com/otp-auth/internal/application/ports/repositories"
	"github.com/otp-auth/internal/application/ports/services"
	"github.com/otp-auth/internal/application/usecases"
	"github.com/otp-auth/internal/config"
	"github.com/otp-auth/internal/infrastructure/http/handlers"
	"github.com/otp-auth/internal/infrastructure/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterConfig holds router configuration
type RouterConfig struct {
	CORSConfig      middleware.CORSConfig
	LoggingConfig   middleware.LoggingConfig
	RateLimitConfig middleware.RateLimitConfig
	EnableSwagger   bool
	TrustedProxies  []string
}

// DefaultRouterConfig returns a default router configuration
func DefaultRouterConfig() RouterConfig {
	return RouterConfig{
		CORSConfig:      middleware.DefaultCORSConfig(),
		LoggingConfig:   middleware.DefaultLoggingConfig(),
		RateLimitConfig: middleware.DefaultRateLimitConfig(),
		EnableSwagger:   true,
		TrustedProxies:  []string{"127.0.0.1", "172.25.0.0/16"},
	}
}

// Dependencies holds all the dependencies needed for the router
type Dependencies struct {
	// Use cases
	SendOTPUseCase        *usecases.SendOTPUseCase
	LoginUseCase          *usecases.LoginUseCase
	RefreshUseCase        *usecases.RefreshUseCase
	LogoutUseCase         *usecases.LogoutUseCase
	GetUserProfileUseCase *usecases.GetUserProfileUseCase
	GetUsersListUseCase   *usecases.GetUsersListUseCase

	// Services
	JWTService services.JWTService

	// Repositories
	RateLimiter repositories.RateLimiter

	// Configuration
	RateLimitConfig *config.RateLimitConfig
}

// SetupRouter sets up the Gin router with all routes and middleware
func SetupRouter(deps Dependencies, config RouterConfig) *gin.Engine {
	router := gin.New()

	// Set trusted proxies
	router.SetTrustedProxies(config.TrustedProxies)

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logging(config.LoggingConfig))
	router.Use(middleware.CORS(config.CORSConfig))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(deps.SendOTPUseCase, deps.LoginUseCase, deps.RefreshUseCase, deps.LogoutUseCase)
	userHandler := handlers.NewUserHandler(deps.GetUserProfileUseCase, deps.GetUsersListUseCase)
	healthHandler := handlers.NewHealthHandler()

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(deps.JWTService)

	// Health check routes (no authentication required)
	router.GET("/health", healthHandler.Health)
	router.GET("/live", healthHandler.Live)
	router.GET("/ready", healthHandler.Ready)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// IP rate limiting /api/v1 endpoint group
		v1.Use(middleware.IPBasedRateLimit(deps.RateLimiter, deps.RateLimitConfig.Requests, deps.RateLimitConfig.Window))
		// Authentication routes (no authentication required)
		auth := v1.Group("/auth")
		{
			// Rate limit for OTP sending (per phone number)
			auth.POST("/send-otp",
				middleware.OtpRateLimit(deps.RateLimiter, deps.RateLimitConfig.OTPLimit, deps.RateLimitConfig.OTPWindow),
				authHandler.SendOTP,
			)

			auth.POST("/login",
				authHandler.Login,
			)

			auth.POST("/refresh",
				authHandler.RefreshToken,
			)

			auth.POST("/logout",
				authHandler.Logout,
			)
		}

		// User routes
		users := v1.Group("/users")
		{
			// Get current user profile (authentication required)
			users.GET("/profile",
				authMiddleware.RequireAuth(),
				userHandler.GetProfile,
			)

			// TODO: Admin routes (admin authentication required)
			users.GET("/",
				authMiddleware.RequireAuth(),
				userHandler.GetUsers,
			)
		}
	}

	// Swagger documentation (if enabled)
	if config.EnableSwagger {
		// Serve the OpenAPI YAML content directly
		router.StaticFile("/openapi.yaml", "./openapi.yaml")

		// Configure Swagger UI to use our custom OpenAPI file
		url := ginSwagger.URL("http://localhost:8080/openapi.yaml") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	// Catch-all route for 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Route not found",
			"code":  "NOT_FOUND",
		})
	})

	return router
}

// SetupProductionRouter sets up a production-ready router
func SetupProductionRouter(deps Dependencies, allowedOrigins []string) *gin.Engine {
	config := RouterConfig{
		CORSConfig:      middleware.ProductionCORSConfig(allowedOrigins),
		LoggingConfig:   middleware.DefaultLoggingConfig(),
		RateLimitConfig: middleware.DefaultRateLimitConfig(),
		EnableSwagger:   false, // Disable swagger in production
		TrustedProxies:  []string{"127.0.0.1", "10.0.0.0/8", "172.25.0.0/16", "192.168.0.0/16"},
	}

	return SetupRouter(deps, config)
}

// SetupDevelopmentRouter sets up a development router with detailed logging
func SetupDevelopmentRouter(deps Dependencies) *gin.Engine {
	config := DefaultRouterConfig()
	config.LoggingConfig.LogBody = true // Enable body logging in development
	config.EnableSwagger = true

	return SetupRouter(deps, config)
}
