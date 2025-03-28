package utils

import (
	"crypto/rand"
	"encoding/base64"
	"final/config"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// generateSalt creates a random salt of specified length
func generateSalt(length int) (string, error) {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}

// HashPassword hashes the password with a custom salt
func HashPassword(password string) (string, error) {
    // Generate salt (16 bytes is common)
    salt, err := generateSalt(16)
    if err != nil {
        return "", err
    }
    
    // Combine password with salt
    saltedPassword := fmt.Sprintf("%s%s", password, salt)
    
    // Hash the salted password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), config.BCryptCost())
    if err != nil {
        return "", err
    }
    
    // Store both the salt and hash to be able to verify later
    // Format: base64(salt):base64(hash)
    return fmt.Sprintf("%s:%s", salt, string(hashedPassword)), nil
}

// CheckPasswordHash verifies if the password matches the stored hash
func CheckPasswordHash(password, storedValue string) bool {
    // Split the stored value to get the salt and hash
    parts := strings.Split(storedValue, ":")
    if len(parts) != 2 {
        return false
    }
    
    salt, storedHash := parts[0], parts[1]
    
    // Recreate the salted password
    saltedPassword := fmt.Sprintf("%s%s", password, salt)
    
    // Compare the hash
    err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(saltedPassword))
    return err == nil
}