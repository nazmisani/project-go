package main

import (
	database "final/config"
	"final/controllers"
	"final/models"
	"final/routes"
	"log"

	"github.com/joho/godotenv"
)

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
