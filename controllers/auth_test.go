package controllers

import (
	"bytes"
	"encoding/json"
	"final/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
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
}

func TestLogin(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/login", Login)

	// Test data
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}

	// Convert user struct to JSON
	jsonValue, _ := json.Marshal(user)

	// Create a request
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	resp := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
}