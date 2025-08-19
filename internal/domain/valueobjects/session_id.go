package valueobjects

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

// SessionID represents a session identifier
type SessionID string

const (
	// SessionIDLength is the length of session ID in bytes
	SessionIDLength = 32
)

// NewSessionID creates a new random session ID
func NewSessionID() (SessionID, error) {
	bytes := make([]byte, SessionIDLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.New("failed to generate session ID")
	}
	return SessionID(hex.EncodeToString(bytes)), nil
}

// NewSessionIDFromString creates a session ID from an existing string
func NewSessionIDFromString(sessionID string) (SessionID, error) {
	// Remove any whitespace
	sessionID = strings.TrimSpace(sessionID)
	
	if sessionID == "" {
		return "", errors.New("session ID cannot be empty")
	}

	// Validate the session ID format (should be hex string)
	if len(sessionID) != SessionIDLength*2 {
		return "", errors.New("invalid session ID length")
	}

	// Check if it's a valid hex string
	if _, err := hex.DecodeString(sessionID); err != nil {
		return "", errors.New("invalid session ID format")
	}

	return SessionID(sessionID), nil
}

// String returns the string representation of the session ID
func (s SessionID) String() string {
	return string(s)
}

// IsValid checks if the session ID is valid
func (s SessionID) IsValid() bool {
	sessionID := string(s)
	if len(sessionID) != SessionIDLength*2 {
		return false
	}
	_, err := hex.DecodeString(sessionID)
	return err == nil
}

// IsEmpty checks if the session ID is empty
func (s SessionID) IsEmpty() bool {
	return string(s) == ""
}