package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/otp-auth/pkg/errors"
)

// Config holds Redis connection configuration
type Config struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig returns a default Redis configuration
func DefaultConfig() Config {
	return Config{
		Addr:         "localhost:6379",
		Password:     "", // No password by default
		DB:           0,  // Default DB
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}
}

// NewConnection creates a new Redis client connection
func NewConnection(config Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, errors.NewInternalError("Failed to connect to Redis", err)
	}

	return client, nil
}

// NewClusterConnection creates a new Redis cluster client connection
func NewClusterConnection(addrs []string, password string) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, errors.NewInternalError("Failed to connect to Redis cluster", err)
	}

	return client, nil
}

// HealthCheck performs a health check on the Redis connection
func HealthCheck(ctx context.Context, client *redis.Client) error {
	if err := client.Ping(ctx).Err(); err != nil {
		return errors.NewInternalError("Redis health check failed", err)
	}
	return nil
}

// ClusterHealthCheck performs a health check on the Redis cluster connection
func ClusterHealthCheck(ctx context.Context, client *redis.ClusterClient) error {
	if err := client.Ping(ctx).Err(); err != nil {
		return errors.NewInternalError("Redis cluster health check failed", err)
	}
	return nil
}