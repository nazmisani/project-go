package utils

import (
	"errors"
	"final/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenDetails struktur untuk menyimpan informasi token
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// GenerateJWT membuat token JWT access + refresh
func GenerateJWT(username string) (*TokenDetails, error) {
	td := &TokenDetails{}
	
	// Set waktu kedaluwarsa
	td.AtExpires = time.Now().Add(time.Hour * time.Duration(config.JWTExpiryTime())).Unix()
	td.RtExpires = time.Now().Add(time.Hour * time.Duration(config.JWTRefreshExpiryTime())).Unix()

	// Access token
	accessClaims := jwt.MapClaims{
		"username": username,
		"exp":      td.AtExpires,
		"type":     "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	var err error
	td.AccessToken, err = accessToken.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"username": username,
		"exp":      td.RtExpires,
		"type":     "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	td.RefreshToken, err = refreshToken.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	return td, nil
}

// ParseJWT memverifikasi dan membaca token JWT
func ParseJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Memastikan signing method adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	
	return nil, errors.New("invalid token")
}

// ParseRefreshToken khusus untuk memvalidasi refresh token
func ParseRefreshToken(tokenString string) (*jwt.MapClaims, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return nil, err
	}
	
	// Verifikasi bahwa ini adalah refresh token
	if tokenType, ok := (*claims)["type"].(string); !ok || tokenType != "refresh" {
		return nil, errors.New("not a refresh token")
	}
	
	return claims, nil
}
