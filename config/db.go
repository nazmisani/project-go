package config

import (
	"final/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Error: DATABASE_URL tidak ditemukan. Silakan atur DATABASE_URL di Railway dashboard atau file .env")
	}

	// Mencoba melakukan koneksi ke database dengan retry
	var db *gorm.DB
	var err error
	
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[error] failed to initialize database, got error %v", err)
		log.Fatal("Gagal terhubung ke database. Pastikan DATABASE_URL sudah benar dan database PostgreSQL sudah berjalan")
	}

	DB = db

	// Migrasi model ke database
	if err := DB.AutoMigrate(&models.User{}, &models.Post{}); err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	log.Println("Berhasil terhubung ke database dan melakukan migrasi")
	// Cek koneksi database
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan koneksi database: %v", err)
	}

	// Test koneksi
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}
}

