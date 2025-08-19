package main

import (
	"context"
	"database/sql"
	"github.com/otp-auth/internal/application/ports/repositories"
	"github.com/otp-auth/internal/application/ports/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	redisClient "github.com/go-redis/redis/v8"

	"github.com/otp-auth/internal/application/usecases"
	"github.com/otp-auth/internal/config"
	"github.com/otp-auth/internal/infrastructure/http/router"
	"github.com/otp-auth/internal/infrastructure/persistence/postgres"
	"github.com/otp-auth/internal/infrastructure/persistence/redis"
	infraServices "github.com/otp-auth/internal/infrastructure/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load("./configs")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connection
	db, err := postgres.NewConnection(postgres.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := postgres.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize Redis connection
	redisConn, err := redis.NewConnection(redis.Config{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		IdleTimeout:  cfg.Redis.IdleTimeout,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	userRepo, otpRepo, tokenRepo, rateLimiter := initializeRepositories(db, redisConn)

	// Initialize services
	otpSender, jwtService, hashService := initializeServices(cfg)

	// Initialize use cases
	sendOTPUseCase := usecases.NewSendOTPUseCase(
		userRepo, otpRepo, rateLimiter,
		otpSender, hashService,
		cfg.OTP.TTL,
		cfg.Security.RateLimit.OTPWindow,
		cfg.Security.RateLimit.OTPLimit,
	)

	loginUseCase := usecases.NewLoginUseCase(
		userRepo, otpRepo, tokenRepo,
		jwtService, hashService,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	refreshUseCase := usecases.NewRefreshUseCase(
		userRepo, tokenRepo,
		jwtService, hashService,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	logoutUseCase := usecases.NewLogoutUseCase(
		tokenRepo,
		hashService,
	)

	getUserProfileUseCase := usecases.NewGetUserProfileUseCase(
		userRepo,
	)

	getUsersListUseCase := usecases.NewGetUsersListUseCase(
		userRepo,
	)

	log.Println("ratelimit: ", cfg.Security.RateLimit.OTPWindow, cfg.Security.RateLimit.OTPLimit)
	// Setup router
	deps := router.Dependencies{
		SendOTPUseCase:        sendOTPUseCase,
		LoginUseCase:          loginUseCase,
		RefreshUseCase:        refreshUseCase,
		LogoutUseCase:         logoutUseCase,
		GetUserProfileUseCase: getUserProfileUseCase,
		GetUsersListUseCase:   getUsersListUseCase,
		JWTService:            jwtService,
		RateLimiter:           rateLimiter,
		RateLimitConfig:       &cfg.Security.RateLimit,
	}

	var r *gin.Engine
	if cfg.Server.Mode == gin.ReleaseMode {
		r = router.SetupProductionRouter(deps, []string{"*"}) // Allow all origins for now
	} else {
		r = router.SetupDevelopmentRouter(deps)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         cfg.Server.GetAddress(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.Server.GetAddress())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	db.Close()

	// Close Redis connection
	redisConn.Close()

	log.Println("Server exited")
}

func initializeRepositories(db *sql.DB, redisConn *redisClient.Client) (repositories.UserRepository, repositories.OTPRepository, repositories.TokenRepository, repositories.RateLimiter) {
	return postgres.NewUserRepository(db), redis.NewOTPRepository(redisConn), postgres.NewTokenRepository(db), redis.NewRateLimiter(redisConn)
}

func initializeServices(cfg *config.Config) (services.OTPSender, services.JWTService, services.HashService) {
	// Initialize JWT service
	jwtService, err := infraServices.NewECDSAJWTService(infraServices.JWTConfig{
		PrivateKeyPEM:   cfg.JWT.PrivateKeyPEM,
		PublicKeyPEM:    cfg.JWT.PublicKeyPEM,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		Issuer:          cfg.JWT.Issuer,
	})
	if err != nil {
		log.Fatalf("Failed to initialize JWT service: %v", err)
	}

	// Initialize hash service
	hashService := infraServices.NewBcryptHashService(infraServices.HashConfig{
		Cost: cfg.Hash.Cost,
	})

	// Initialize OTP sender
	var otpSender services.OTPSender
	switch cfg.OTP.SenderType {
	case "console":
		otpSender = infraServices.NewConsoleOTPSender(nil)
	default:
		// Default to console sender
		otpSender = infraServices.NewConsoleOTPSender(nil)
		log.Printf("Unknown OTP sender type '%s', using console sender", cfg.OTP.SenderType)
	}

	return otpSender, jwtService, hashService
}
