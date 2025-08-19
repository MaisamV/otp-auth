package services

// HashService defines the interface for hashing operations
type HashService interface {
	// Hash hashes a plain text string
	Hash(plaintext string) (string, error)
	
	// Compare compares a plain text string with a hash
	Compare(plaintext, hash string) bool
	
	// HashPassword hashes a password using bcrypt
	HashPassword(password string) (string, error)
	
	// VerifyPassword verifies a password against its hash
	VerifyPassword(password, hash string) error
	
	// HashOTP hashes an OTP code using bcrypt
	HashOTP(otp string) (string, error)
	
	// VerifyOTP verifies an OTP code against its hash
	VerifyOTP(otp, hash string) error
	
	// HashRefreshToken hashes a refresh token using bcrypt
	HashRefreshToken(token string) (string, error)
	
	// VerifyRefreshToken verifies a refresh token against its hash
	VerifyRefreshToken(token, hash string) error
	
	// GenerateRandomString generates a random string of specified length
	GenerateRandomString(length int) (string, error)
}