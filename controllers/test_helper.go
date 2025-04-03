package controllers

import (
	"final/config"
	"final/models"
	"final/utils"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var originalDB *gorm.DB

// SetupTestDB initializes a test database connection
func SetupTestDB() *gorm.DB {
	// Load .env file if exists
	_ = godotenv.Load()

	// Use the same database as the application but with transaction
	dsn := "host=localhost user=postgres password=postgres dbname=testdb port=5432 sslmode=disable"
	
	// Override with environment variables if available
	if os.Getenv("TEST_DB_DSN") != "" {
		dsn = os.Getenv("TEST_DB_DSN")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Store the original DB connection
	originalDB = config.DB
	
	// Start a transaction that will be rolled back
	tx := db.Begin()

	// Set the global DB variable for controllers to use the transaction
	config.DB = tx
	testDB = tx

	return tx
}

// ClearTestDB rolls back the test transaction and restores the original DB connection
func ClearTestDB() {
	if testDB != nil {
		// Rollback the transaction
		testDB.Rollback()
		
		// Restore the original DB connection
		if originalDB != nil {
			config.DB = originalDB
		}
	}
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T) models.User {
	// Hash the password
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Create a test user
	user := models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: hashedPassword,
	}

	// Save to database
	result := testDB.Create(&user)
	if result.Error != nil {
		t.Fatalf("Failed to create test user: %v", result.Error)
	}

	return user
}

// SetupTestRouter returns a configured Gin router for testing
func SetupTestRouter() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create a new router
	r := gin.Default()
	
	return r
}

// TestMain is used for setup and teardown of tests
func TestMain(m *testing.M) {
	// Setup
	SetupTestDB()
	
	// Run tests
	exitCode := m.Run()
	
	// Cleanup
	ClearTestDB()
	
	// Exit
	os.Exit(exitCode)
}

// RunWithTransaction runs a test function within a transaction that will be rolled back
func RunWithTransaction(t *testing.T, testFunc func(t *testing.T)) {
	// Setup test database with transaction
	SetupTestDB()
	
	// Run the test
	testFunc(t)
	
	// Cleanup - rollback transaction
	ClearTestDB()
}