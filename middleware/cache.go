package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sort"
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
	// Jumlah maksimum item dalam cache
	maxCacheItems := 1000
	// Interval pembersihan cache (5 menit)
	cleanupInterval := 5 * time.Minute
	
	// Goroutine untuk membersihkan cache yang sudah kedaluwarsa
	go func() {
		for {
			time.Sleep(cleanupInterval)
			now := time.Now()
			mutex.Lock()
			
			// Hapus item yang sudah kedaluwarsa
			for k, v := range cache {
				if now.After(v.Expiration) {
					delete(cache, k)
				}
			}
			
			// Jika cache masih terlalu besar, hapus 20% item tertua
			if len(cache) > maxCacheItems {
				// Buat slice untuk menyimpan key dan waktu kedaluwarsa
				items := make([]struct {
					key string
					exp time.Time
				}, 0, len(cache))
				
				for k, v := range cache {
					items = append(items, struct {
						key string
						exp time.Time
					}{k, v.Expiration})
				}
				
				// Urutkan berdasarkan waktu kedaluwarsa (tertua dulu)
				sort.Slice(items, func(i, j int) bool {
					return items[i].exp.Before(items[j].exp)
				})
				
				// Hapus 20% item tertua
				toRemove := len(cache) / 5
				for i := 0; i < toRemove; i++ {
					delete(cache, items[i].key)
				}
			}
			
			mutex.Unlock()
		}
	}()

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

// ClearCache untuk menghapus semua cache
func ClearCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Cache berhasil dihapus",
		})
	}
}