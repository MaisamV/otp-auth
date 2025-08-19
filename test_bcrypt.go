package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Test the hash from logs
	hash := "$2a$10$apuXN5C11Ufzlr3DzEPPX.W6VGurDYqXTzfcAedhaKsI5cvIbeAWq"
	otp := "155204"
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(otp))
	if err != nil {
		fmt.Printf("Verification failed: %v\n", err)
	} else {
		fmt.Println("Verification successful!")
	}
	
	// Also test generating a new hash for the OTP
	newHash, err := bcrypt.GenerateFromPassword([]byte(otp), 10)
	if err != nil {
		fmt.Printf("Hash generation failed: %v\n", err)
	} else {
		fmt.Printf("New hash for %s: %s\n", otp, string(newHash))
		
		// Verify the new hash
		err = bcrypt.CompareHashAndPassword(newHash, []byte(otp))
		if err != nil {
			fmt.Printf("New hash verification failed: %v\n", err)
		} else {
			fmt.Println("New hash verification successful!")
		}
	}
}