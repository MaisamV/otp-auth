package entities

import (
	"time"

	"github.com/otp-auth/internal/domain/valueobjects"
)

// User represents a user in the system
type User struct {
	ID          string                      `json:"id"`
	PhoneNumber valueobjects.PhoneNumber    `json:"phone_number"`
	Scope       string                      `json:"scope"` // empty for normal users, "superadmin" for admin access
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

// NewUser creates a new user with the given phone number
func NewUser(phoneNumber valueobjects.PhoneNumber) *User {
	now := time.Now()
	return &User{
		PhoneNumber: phoneNumber,
		Scope:       "user", // Default to normal user
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Scope == "superadmin"
}

// UpdateScope updates the user's scope
func (u *User) UpdateScope(scope string) {
	u.Scope = scope
	u.UpdatedAt = time.Now()
}