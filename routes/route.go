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

	// Auth Routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.Use(middleware.AuthMiddleware())
	// User Routes
	r.POST("/users", controllers.CreateUser)
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUser)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	// Post Routes
	r.POST("/posts", controllers.CreatePost)

	
	r.GET("/users/post", controllers.GetUsersWithPosts)

	return r
}
