package controllers

import (
	"bytes"
	"encoding/json"
	"final/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegister tests the Register endpoint
func TestRegister(t *testing.T) {
	RunWithTransaction(t, func(t *testing.T) {
		// Setup router
		r := SetupTestRouter()
		r.POST("/register", Register)

		// Test data
		user := models.User{
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "password123",
		}

		// Convert user struct to JSON
		jsonValue, _ := json.Marshal(user)

		// Create a request
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		resp := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Contains(t, resp.Body.String(), "User registered successfully")
		
		// Verify user was created in database
		var createdUser models.User
		result := testDB.Where("username = ?", user.Username).First(&createdUser)
		assert.Nil(t, result.Error)
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.Email, createdUser.Email)
	})
}

func TestLogin(t *testing.T) {
	RunWithTransaction(t, func(t *testing.T) {
		// Create a test user in the database
		testUser := CreateTestUser(t)
		
		// Setup router
		r := SetupTestRouter()
		r.POST("/login", Login)

		// Test data - use the same username as the created user
		loginData := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: testUser.Username,
			Password: "password123", // Plain password before hashing
		}

		// Convert login data to JSON
		jsonValue, _ := json.Marshal(loginData)

		// Create a request
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		resp := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "access_token")
		assert.Contains(t, resp.Body.String(), "refresh_token")
	})
}