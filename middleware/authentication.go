package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"final/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware untuk validasi JWT dengan caching
func AuthMiddleware() gin.HandlerFunc {
	// Cache untuk menyimpan token yang valid (token -> username)
	tokenCache := make(map[string]string)
	// Mutex untuk mengamankan akses ke cache
	var mutex sync.Mutex
	
	return func(c *gin.Context) {
		// Ambil token dari header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Authorization header diperlukan",
			})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		
		// Cek cache terlebih dahulu
		mutex.Lock()
		username, exists := tokenCache[tokenString]
		mutex.Unlock()
		
		if exists {
			// Token ada di cache, set username ke context
			c.Set("username", username)
			c.Next()
			return
		}

		// Token tidak ada di cache, validasi dengan JWT
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Token tidak valid",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Verifikasi bahwa token belum kedaluwarsa
		if exp, ok := (*claims)["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
		}

		// Verifikasi bahwa ini adalah access token bukan refresh token
		if tokenType, ok := (*claims)["type"].(string); ok && tokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		// Set username untuk digunakan dalam handler
		username, ok := (*claims)["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		c.Set("username", username)
		
		// Simpan token di cache
		mutex.Lock()
		tokenCache[tokenString] = username
		mutex.Unlock()
		
		c.Next()
	}
}
