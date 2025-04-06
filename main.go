package main

import (
	database "final/config"
	"final/controllers"
	"final/models"
	"final/routes"
	"log"

	_ "final/docs" // Import docs untuk Swagger

	"github.com/joho/godotenv"
)

// @title           Final Project API
// @version         1.0
// @description     API untuk Final Project dengan fitur user dan post management
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Inisialisasi Cloudinary
	controllers.InitCloudinary()

	// Koneksi ke database
	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.User{} ,&models.Post{}) // Migrate database

	// Setup router
	r := routes.SetupRouter()
	r.Run(":8080") // Jalankan server di port 8080
}
