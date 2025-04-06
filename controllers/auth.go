package controllers

import (
	"final/config"
	"final/models"
	"final/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body docs.RegisterRequest true "User registration data"
// @Success 201 {object} docs.UserResponse "User created successfully"
// @Failure 400 {object} docs.ErrorResponse "Bad request - validation error or username/email already exists"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /register [post]
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Validasi gagal: " + err.Error(),
		})
		return
	}

	// Validasi tambahan untuk username
	if len(input.Username) < 3 || len(input.Username) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Username harus antara 3-50 karakter",
		})
		return
	}

	// Validasi bahwa username/email belum digunakan
	var existingUser models.User
	if result := config.DB.Where("username = ?", input.Username).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Username sudah digunakan",
		})
		return
	}

	if result := config.DB.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Email sudah digunakan",
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error hashing password: " + err.Error(),
		})
		return
	}

	// Simpan user ke database
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     "user", // Set default role
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Gagal membuat user: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "User berhasil didaftarkan",
		"data": gin.H{
			"userId":   user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body docs.LoginRequest true "User login credentials"
// @Success 200 {object} docs.TokenResponse "Login successful"
// @Failure 400 {object} docs.ErrorResponse "Bad request - validation error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid credentials"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /login [post]
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari user di database
	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Cek password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token JWT (access + refresh)
	tokens, err := utils.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    config.JWTExpiryTime() * 3600, // Dalam detik
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token" example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
// @Success 200 {object} docs.TokenResponse "New tokens generated successfully"
// @Failure 400 {object} docs.ErrorResponse "Bad request - invalid refresh token"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - expired or invalid token"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /refresh [post]
func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dan validasi refresh token
	claims, err := utils.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Ambil username dari claims
	username, ok := (*claims)["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Verifikasi user masih ada di database
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate token baru
	tokens, err := utils.GenerateJWT(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    config.JWTExpiryTime() * 3600, // Dalam detik
	})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate user's refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid token"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /logout [post]
func Logout(c *gin.Context) {
	// Di implementasi lengkap, kita akan menambahkan token ke blacklist
	// Tapi untuk sekarang, kita hanya kembalikan sukses
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
