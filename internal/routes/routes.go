package routes

import (
	"gin-quickstart/internal/config"
	"gin-quickstart/internal/handlers"
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	apiLogRepo repositories.APILogRepository,
) {
	// Global Middleware
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())
	r.Use(gin.Recovery())

	// Rate limiter: 100 requests per minute
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	r.Use(rateLimiter.Middleware())

	// API Monitor
	r.Use(middleware.APIMonitorMiddleware(apiLogRepo))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/google", authHandler.GoogleLogin)
			auth.GET("/google/callback", authHandler.GoogleCallback)
			auth.GET("/github", authHandler.GithubLogin)
			auth.GET("/github/callback", authHandler.GithubCallback)
		}

		// User routes (protected)
		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware(cfg))
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
			user.DELETE("/account", userHandler.DeleteAccount)
		}

		// Admin routes (protected + admin only)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.Use(middleware.AdminMiddleware())
		{
			// Users management
			admin.GET("/users", userHandler.GetAllUsers)
			admin.GET("/users/:id", userHandler.GetUserByID)
			admin.PUT("/users/:id", userHandler.UpdateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUser)

			// Dashboard & analytics
			// TODO: Add dashboard handlers
		}
	}
}
