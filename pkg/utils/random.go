package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// GenerateRandomBytes generates random bytes of specified length
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateRandomNumber generates a random number between min and max (inclusive)
func GenerateRandomNumber(min, max int64) (int64, error) {
	if min >= max {
		return min, nil
	}
	
	n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return 0, err
	}
	
	return n.Int64() + min, nil
}

// GenerateOTP generates a numeric OTP of specified length
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	
	min := int64(1)
	max := int64(1)
	for i := 0; i < length; i++ {
		if i == 0 {
			min = 1
			max = 9
		} else {
			min *= 10
			max = max*10 + 9
		}
	}
	
	num, err := GenerateRandomNumber(min, max)
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%0*d", length, num), nil
}