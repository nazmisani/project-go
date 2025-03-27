package config

import (
	"os"
)

// GetJWTSecret mengembalikan secret key untuk JWT
// Mengambil dari environment variable atau menggunakan default jika tidak ada
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback ke default key (seharusnya tidak digunakan di production)
		secret = "secret_key_for_development_only_change_in_production"
	}
	return []byte(secret)
}

// JWTExpiryTime adalah waktu kedaluwarsa token dalam jam
func JWTExpiryTime() int {
	return 24 // 24 jam = 1 hari
}

// JWTRefreshExpiryTime adalah waktu kedaluwarsa refresh token dalam jam
func JWTRefreshExpiryTime() int {
	return 24 * 7 // 7 hari
}

// BCryptCost adalah cost untuk hashing bcrypt
// Minimal 12 untuk keamanan yang baik
func BCryptCost() int {
	return 12
} 