package main

import (
	"log"
	"os"

	"otp-auth-service/docs"
	"otp-auth-service/internal/config"
	"otp-auth-service/internal/handlers"
	"otp-auth-service/internal/middleware"
	"otp-auth-service/internal/repository"
	"otp-auth-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title OTP Authentication Service
// @version 1.0
// @description A backend service for OTP-based authentication and user management
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	otpRepo := repository.NewOTPRepository(redisClient)

	// Initialize services
	authService := services.NewAuthService(userRepo, otpRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Setup Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/request-otp", authHandler.RequestOTP)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
		}

		// User routes (protected)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			users.GET("/", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
		}
	}

	// Swagger documentation
	docs.SwaggerInfo.Title = "OTP Authentication Service"
	docs.SwaggerInfo.Description = "A backend service for OTP-based authentication and user management"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
