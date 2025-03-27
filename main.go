package main

import (
	database "final/config"
	"final/models"
	"final/routes"
)

func main() {
	// Koneksi ke database
	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.User{} ,&models.Post{}) // Migrate database

	// Setup router
	r := routes.SetupRouter()
	r.Run(":8080") // Jalankan server di port 8080
}
