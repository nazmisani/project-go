package controllers

import (
	"context"
	"final/config"
	"final/models"
	"log"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&post)
	c.JSON(http.StatusOK, post)
}
// Konfigurasi Cloudinary
var cld, err = cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
func init() {
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
}	

// UploadToCloudinary handles file uploads to Cloudinary
func UploadToCloudinary(c *gin.Context) {
	// Ambil file dari request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mendapatkan file"})
		return
	}

	// Buka file untuk upload
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer src.Close()

	// Upload ke Cloudinary
	uploadResult, err := cld.Upload.Upload(context.Background(), src, uploader.UploadParams{})
	if err != nil {
		log.Println("Upload error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload ke Cloudinary"})
		return
	}

	// Beri respon ke client
	c.JSON(http.StatusOK, gin.H{
		"message": "File berhasil diupload!",
		"url":     uploadResult.SecureURL,
	})
}
