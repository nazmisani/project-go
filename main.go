package main

import (
	database "final/config"
	"final/controllers"
	"final/routes"
	"log"
	"os"

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
		log.Println("Info: .env file tidak ditemukan. Pastikan environment variables sudah diatur di Railway")
	}

	// Periksa DATABASE_URL
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("Error: DATABASE_URL tidak ditemukan. Silakan atur DATABASE_URL di Railway dashboard")
	}

	// Periksa CLOUDINARY_URL
	if os.Getenv("CLOUDINARY_URL") == "" {
		log.Fatal("Error: CLOUDINARY_URL tidak ditemukan. Silakan atur CLOUDINARY_URL di Railway dashboard")
	}

	// Inisialisasi Cloudinary
	controllers.InitCloudinary()

	// Koneksi ke database
	database.ConnectDatabase()
	// Migrasi database sudah dilakukan di ConnectDatabase()

	// Setup router
	r := routes.SetupRouter()

	// Jalankan server di port 8080
	log.Println("Server berjalan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
