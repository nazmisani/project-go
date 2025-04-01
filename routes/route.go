package routes

import (
	"final/controllers"
	"final/middleware"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		startTime := time.Now() // Merekam waktu mulai eksekusi request
		context.Next()          // Melanjutkan ke handler berikutnya
		duration := time.Since(startTime) // Menghitung durasi eksekusi

		log.Printf("%s %s | Status: %d | Duration: %s",
			context.Request.Method,
			context.Request.URL.Path,
			context.Writer.Status(),
			duration)
	}
}

func RateLimitMiddleware(rateLimiter *rate.Limiter) gin.HandlerFunc {
	return func(context *gin.Context) {
		if !rateLimiter.Allow() { // Mengecek apakah request diperbolehkan berdasarkan rate limiter
			context.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			context.Abort() // Menghentikan eksekusi request
			return
		}
		context.Next() // Melanjutkan request jika diperbolehkan
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(LoggingMiddleware())

	// Rate Limiter
	requestLimiter := rate.NewLimiter(1, 5)
	r.Use(RateLimitMiddleware(requestLimiter))

	// Public Auth Routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.RefreshToken) // Endpoint untuk refresh token
	
	// Protected Routes (require valid JWT)
	authRoutes := r.Group("/")
	authRoutes.Use(middleware.AuthMiddleware())
	
	// Logout endpoint
	authRoutes.POST("/logout", controllers.Logout)
	
	// User Routes
	authRoutes.POST("/users", controllers.CreateUser)
	authRoutes.GET("/users", controllers.GetUsers)
	authRoutes.GET("/users/:id", controllers.GetUser)
	authRoutes.PUT("/users/:id", controllers.UpdateUser)
	authRoutes.DELETE("/users/:id", controllers.DeleteUser)

	// Post Routes
	authRoutes.POST("/posts", controllers.CreatePost)
	authRoutes.GET("/users/post", controllers.GetUsersWithPosts)

	// Upload Route
	authRoutes.POST("/upload", controllers.UploadToCloudinary)

	return r
}
