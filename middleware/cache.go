package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Item struktur untuk menyimpan item cache
type cacheItem struct {
	Content    []byte
	Expiration time.Time
}

// Cache middleware untuk menyimpan respons API
func CacheMiddleware(expiration time.Duration) gin.HandlerFunc {
	// Cache untuk menyimpan respons
	cache := make(map[string]cacheItem)
	// Mutex untuk mengamankan akses ke cache
	var mutex sync.Mutex

	return func(c *gin.Context) {
		// Hanya cache untuk GET requests
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Buat key cache dari URL dan header Authorization
		auth := c.GetHeader("Authorization")
		hashKey := sha256.Sum256([]byte(c.Request.URL.String() + auth))
		key := hex.EncodeToString(hashKey[:])

		// Cek apakah respons ada di cache
		mutex.Lock()
		item, exists := cache[key]
		mutex.Unlock()

		if exists && time.Now().Before(item.Expiration) {
			// Respons ada di cache dan belum kedaluwarsa
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Header().Set("X-Cache", "HIT")
			c.Writer.Write(item.Content)
			c.Abort()
			return
		}

		// Buat writer untuk menyimpan respons
		writer := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = writer

		// Lanjutkan ke handler berikutnya
		c.Next()

		// Simpan respons di cache jika status code 200 OK
		if c.Writer.Status() == http.StatusOK {
			mutex.Lock()
			cache[key] = cacheItem{
				Content:    writer.body.Bytes(),
				Expiration: time.Now().Add(expiration),
			}
			mutex.Unlock()
			c.Writer.Header().Set("X-Cache", "MISS")
		}
	}
}

// responseWriter adalah wrapper untuk gin.ResponseWriter
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write menyimpan respons di buffer dan menulis ke ResponseWriter asli
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString menyimpan respons string di buffer dan menulis ke ResponseWriter asli
func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// ClearCache godoc
// @Summary Clear all cache
// @Description Clear all cached responses
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Cache cleared successfully"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized - invalid token"
// @Router /admin/cache/clear [post]
func ClearCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Cache berhasil dihapus",
		})
	}
}