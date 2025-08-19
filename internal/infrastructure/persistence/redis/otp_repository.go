package redis

import (
	"context"
	"fmt"
	"github.com/otp-auth/internal/application/ports/repositories"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
)

// OTPRepository implements the OTP repository using Redis
type OTPRepository struct {
	client *redis.Client
}

// NewOTPRepository creates a new Redis OTP repository
func NewOTPRepository(client *redis.Client) repositories.OTPRepository {
	return &OTPRepository{
		client: client,
	}
}

// GetRedisKey returns the Redis key for storing this OTP
func GetRedisKey(phoneNumber string) string {
	return "otp:" + phoneNumber
}

// GetRedisValue returns the Redis value for storing this OTP
func GetRedisValue(sessionID string, hashedCode string) string {
	return string(sessionID) + "-" + hashedCode
}

// Store stores an OTP in Redis with TTL
func (r *OTPRepository) Store(ctx context.Context, otp *entities.OTP, ttl time.Duration) error {
	key := GetRedisKey(otp.PhoneNumber.String())
	value := GetRedisValue(otp.SessionID.String(), otp.HashedCode)

	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return errors.NewInternalError("Failed to store OTP", err)
	}

	return nil
}

// Get retrieves an OTP by phone number
func (r *OTPRepository) Get(ctx context.Context, phoneNumber valueobjects.PhoneNumber) (*entities.OTP, error) {
	key := GetRedisKey(phoneNumber.String())

	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.NewNotFoundError("OTP not found or expired", nil)
		}
		return nil, errors.NewInternalError("Failed to retrieve OTP", err)
	}

	// Parse value in format "<session_id>-<hashedOtp>"
	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return nil, errors.NewInternalError("Invalid OTP data format", nil)
	}

	sessionID, err := valueobjects.NewSessionIDFromString(parts[0])
	if err != nil {
		return nil, errors.NewInternalError("Invalid session ID in value", err)
	}

	hashedCode := parts[1]

	otp := &entities.OTP{
		PhoneNumber: phoneNumber,
		SessionID:   sessionID,
		HashedCode:  hashedCode,
		// Note: CreatedAt and ExpiresAt will be zero values since we don't store them anymore
	}

	return otp, nil
}

// Exists checks if an OTP exists for the given phone number
func (r *OTPRepository) Exists(ctx context.Context, phoneNumber valueobjects.PhoneNumber) (bool, error) {
	pattern := fmt.Sprintf("otp:%s:*", phoneNumber.String())
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return false, errors.NewInternalError("Failed to check OTP existence", err)
	}

	return len(keys) > 0, nil
}

// Delete removes an OTP from Redis by phone number
func (r *OTPRepository) Delete(ctx context.Context, phoneNumber valueobjects.PhoneNumber) error {
	key := GetRedisKey(phoneNumber.String())
	// Delete all keys for this phone number
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return errors.NewInternalError("Failed to delete OTP", err)
	}

	return nil
}

// GetByPhoneAndSession retrieves an OTP by phone number and session ID (helper method)
func (r *OTPRepository) GetByPhoneAndSession(ctx context.Context, phoneNumber valueobjects.PhoneNumber, sessionID valueobjects.SessionID) (*entities.OTP, error) {
	key := fmt.Sprintf("otp:%s:%s", phoneNumber.String(), sessionID.String())

	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.NewNotFoundError("OTP not found or expired", nil)
		}
		return nil, errors.NewInternalError("Failed to retrieve OTP", err)
	}

	// Parse the stored value (format: "hashedCode:createdAt:expiresAt")
	var hashedCode string
	var createdAt, expiresAt int64

	n, err := fmt.Sscanf(value, "%s:%d:%d", &hashedCode, &createdAt, &expiresAt)
	if err != nil || n != 3 {
		return nil, errors.NewInternalError("Invalid OTP data format", err)
	}

	otp := &entities.OTP{
		PhoneNumber: phoneNumber,
		SessionID:   sessionID,
		HashedCode:  hashedCode,
		CreatedAt:   time.Unix(createdAt, 0),
		ExpiresAt:   time.Unix(expiresAt, 0),
	}

	return otp, nil
}

// RateLimiter implements the rate limiter using Redis
type RateLimiter struct {
	client *redis.Client
}

// NewRateLimiter creates a new Redis rate limiter
func NewRateLimiter(client *redis.Client) repositories.RateLimiter {
	return &RateLimiter{
		client: client,
	}
}

// CheckRateLimit checks if the rate limit is exceeded
func (r *RateLimiter) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	count, err := r.client.Get(ctx, redisKey).Int()
	if err != nil {
		if err == redis.Nil {
			// Key doesn't exist, rate limit not exceeded
			return true, nil
		}
		return false, errors.NewInternalError("Failed to check rate limit", err)
	}

	return count < limit, nil
}

// IncrementCount increments the rate limit counter
func (r *RateLimiter) IncrementCount(ctx context.Context, key string, window time.Duration) error {
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// Use pipeline for atomic operations
	pipe := r.client.Pipeline()
	pipe.Incr(ctx, redisKey)
	pipe.Expire(ctx, redisKey, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errors.NewInternalError("Failed to increment rate limit counter", err)
	}

	return nil
}

// GetCount gets the current rate limit count
func (r *RateLimiter) GetCount(ctx context.Context, key string) (int, error) {
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	count, err := r.client.Get(ctx, redisKey).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, errors.NewInternalError("Failed to get rate limit count", err)
	}

	return count, nil
}

// ResetCount resets the rate limit counter
func (r *RateLimiter) ResetCount(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	err := r.client.Del(ctx, redisKey).Err()
	if err != nil {
		return errors.NewInternalError("Failed to reset rate limit counter", err)
	}

	return nil
}

// CheckAndIncrement checks the current count and increments if under limit
func (r *RateLimiter) CheckAndIncrement(ctx context.Context, key string, limit int, window time.Duration) (bool, int, error) {
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// Get current count
	count, err := r.client.Get(ctx, redisKey).Int()
	if err != nil && err != redis.Nil {
		return false, 0, errors.NewInternalError("Failed to check rate limit", err)
	}

	// If key doesn't exist, count is 0
	if err == redis.Nil {
		count = 0
	}

	// Check if under limit
	if count >= limit {
		return false, count, nil
	}

	// Increment counter atomically
	pipe := r.client.Pipeline()
	pipe.Incr(ctx, redisKey)
	pipe.Expire(ctx, redisKey, window)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, count, errors.NewInternalError("Failed to increment rate limit counter", err)
	}

	return true, count + 1, nil
}

// Reset resets the count for a key (alias for ResetCount to match interface)
func (r *RateLimiter) Reset(ctx context.Context, key string) error {
	return r.ResetCount(ctx, key)
}
