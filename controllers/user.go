package controllers

import (
	database "final/config"
	"final/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with provided data
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.User true "User data"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} docs.ErrorResponse "Bad request - validation error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid token"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan user ke database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.User "List of users"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid token"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.User "User details"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid token"
// @Failure 404 {object} docs.ErrorResponse "User not found"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User "User updated successfully"
// @Failure 400 {object} docs.ErrorResponse "Bad request - validation error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid token"
// @Failure 404 {object} docs.ErrorResponse "User not found"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Cek apakah user ada
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Validasi input JSON
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data user
	if err := database.DB.Model(&user).Updates(updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by user ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid token"
// @Failure 404 {object} docs.ErrorResponse "User not found"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Cek apakah user ada
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Hapus user dari database
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// GetUsersWithPosts godoc
// @Summary Get all users with their posts
// @Description Get a list of all users including their posts
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of users with their posts"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid token"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /users/post [get]
func GetUsersWithPosts(c *gin.Context) {
	var users []models.User

	// Query dengan preloading posts
	if err := database.DB.Preload("Posts").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to fetch users with posts",
		})
		return
	}

	// Jika berhasil, kirim data dengan status code
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    users,
	})
}
