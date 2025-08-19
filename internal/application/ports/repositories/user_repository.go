package repositories

import (
	"context"

	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
)

// UserReader defines read operations for users
type UserReader interface {
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*entities.User, error)
	
	// GetByPhoneNumber retrieves a user by phone number
	GetByPhoneNumber(ctx context.Context, phoneNumber valueobjects.PhoneNumber) (*entities.User, error)
	
	// List retrieves users with pagination and optional search
	List(ctx context.Context, offset, limit int, searchPhone string, searchDateFrom, searchDateTo string) ([]*entities.User, int64, error)
	
	// Exists checks if a user exists by phone number
	Exists(ctx context.Context, phoneNumber valueobjects.PhoneNumber) (bool, error)
}

// UserWriter defines write operations for users
type UserWriter interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error
	
	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error
	
	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error
}

// UserRepository combines read and write operations
type UserRepository interface {
	UserReader
	UserWriter
}