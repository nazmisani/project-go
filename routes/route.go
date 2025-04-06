package routes

import (
	"final/controllers"
	"final/middleware"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"

	_ "final/docs" // Import docs untuk Swagger
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

	// Middleware global
	r.Use(LoggingMiddleware())
	
	// Rate Limiter - meningkatkan batas rate
	requestLimiter := rate.NewLimiter(5, 10) // 5 request per detik dengan burst 10
	r.Use(RateLimitMiddleware(requestLimiter))
	
	// Cache middleware untuk endpoint GET (30 detik)
	r.Use(middleware.CacheMiddleware(30 * time.Second))
	
	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	authRoutes.GET("/posts", controllers.GetPosts)
	authRoutes.GET("/posts/:id", controllers.GetPost)
	authRoutes.PUT("/posts/:id", controllers.UpdatePost)
	authRoutes.DELETE("/posts/:id", controllers.DeletePost)
	authRoutes.GET("/users/post", controllers.GetUsersWithPosts)

	// Upload Route
	authRoutes.POST("/upload", controllers.UploadToCloudinary)
	
	// Cache management route (admin only) - tidak ditampilkan di Swagger
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	// Gunakan endpoint ini tanpa dokumentasi Swagger
	adminRoutes.POST("/cache/clear", middleware.ClearCache())

	return r
}
