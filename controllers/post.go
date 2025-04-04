package controllers

import (
	"context"
	"errors"
	"final/config"
	"final/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
		return
	}
	
	// Validasi data post
	if post.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Judul post tidak boleh kosong",
		})
		return
	}
	
	if post.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Isi post tidak boleh kosong",
		})
		return
	}
	
	if post.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "UserID tidak boleh kosong",
		})
		return
	}
	
	// Cek apakah user ada
	var user models.User
	if err := config.DB.First(&user, post.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "User tidak ditemukan",
		})
		return
	}
	
	// Simpan post ke database
	result := config.DB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menyimpan post",
			"error":   result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Post berhasil dibuat",
		"data":    post,
	})
}
// Konfigurasi Cloudinary
var cld *cloudinary.Cloudinary

// InitCloudinary initializes the Cloudinary client
func InitCloudinary() {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Println("Warning: CLOUDINARY_URL environment variable is not set")
		return
	}
	
	var err error
	cld, err = cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Println("Failed to initialize Cloudinary:", err)
	} else {
		log.Println("Cloudinary initialized successfully")
	}
}	

// UploadToCloudinary handles file uploads to Cloudinary
func UploadToCloudinary(c *gin.Context) {
	// Ambil file dari request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Gagal mendapatkan file",
			"error":   err.Error(),
		})
		return
	}

	// Validasi ukuran file (maksimal 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Ukuran file terlalu besar (maksimal 10MB)",
		})
		return
	}

	// Validasi tipe file (hanya gambar)
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Tipe file tidak didukung (hanya gambar)",
		})
		return
	}

	// Buka file untuk upload
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membuka file",
			"error":   err.Error(),
		})
		return
	}
	defer src.Close()

	// Cek apakah Cloudinary sudah diinisialisasi
	if cld == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Cloudinary belum diinisialisasi",
		})
		return
	}

	// Upload ke Cloudinary dengan parameter tambahan
	uploadParams := uploader.UploadParams{
		Folder:         "uploads",
		ResourceType:   "image",
		Transformation: "q_auto:good", // Kompresi otomatis dengan kualitas baik
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), src, uploadParams)
	if err != nil {
		log.Println("Upload error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal upload ke Cloudinary",
			"error":   err.Error(),
		})
		return
	}

	// Beri respon ke client
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "File berhasil diupload!",
		"data": gin.H{
			"url":       uploadResult.SecureURL,
			"public_id": uploadResult.PublicID,
			"format":    uploadResult.Format,
			"width":     uploadResult.Width,
			"height":    uploadResult.Height,
			"size":      uploadResult.Bytes,
		},
	})
}

// GetPosts mengambil semua post dengan pagination
func GetPosts(c *gin.Context) {
	var posts []models.Post
	
	// Parameter pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	// Hitung total post untuk pagination
	var total int64
	config.DB.Model(&models.Post{}).Count(&total)
	
	// Query dengan pagination dan preload user
	result := config.DB.Preload("User").Limit(limit).Offset(offset).Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mengambil data post",
			"error":   result.Error.Error(),
		})
		return
	}
	
	// Hitung total halaman
	totalPages := (int(total) + limit - 1) / limit
	
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Berhasil mengambil data post",
		"data":    posts,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetPost mengambil post berdasarkan ID
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	
	// Query dengan preload user
	result := config.DB.Preload("User").First(&post, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Post tidak ditemukan",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mengambil data post",
			"error":   result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Berhasil mengambil data post",
		"data":    post,
	})
}

// UpdatePost memperbarui post berdasarkan ID
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	
	// Cek apakah post ada
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Post tidak ditemukan",
		})
		return
	}
	
	// Validasi input JSON
	var input struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
		return
	}
	
	// Validasi data post
	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Judul post tidak boleh kosong",
		})
		return
	}
	
	if input.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Isi post tidak boleh kosong",
		})
		return
	}
	
	// Update post
	updates := map[string]interface{}{
		"title": input.Title,
		"body":  input.Body,
	}
	
	result := config.DB.Model(&post).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal memperbarui post",
			"error":   result.Error.Error(),
		})
		return
	}
	
	// Ambil post yang sudah diupdate
	config.DB.First(&post, id)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Post berhasil diperbarui",
		"data":    post,
	})
}

// DeletePost menghapus post berdasarkan ID
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	
	// Cek apakah post ada
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Post tidak ditemukan",
		})
		return
	}
	
	// Hapus post
	result := config.DB.Delete(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menghapus post",
			"error":   result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Post berhasil dihapus",
	})
}
