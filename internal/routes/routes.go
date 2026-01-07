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
	categoryHandler *handlers.CategoryHandler,
	productHandler *handlers.ProductHandler,
	cartHandler *handlers.CartHandler,
	wishlistHandler *handlers.WishlistHandler,
	orderHandler *handlers.OrderHandler,
	downloadHandler *handlers.DownloadHandler,
	reviewHandler *handlers.ReviewHandler,
	customOrderHandler *handlers.CustomOrderHandler,
	notificationHandler *handlers.NotificationHandler,
	analyticsHandler *handlers.AnalyticsHandler,
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

		// Category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)

			// Admin only
			categoriesAdmin := categories.Group("")
			categoriesAdmin.Use(middleware.AuthMiddleware(cfg))
			categoriesAdmin.Use(middleware.AdminMiddleware())
			{
				categoriesAdmin.POST("", categoryHandler.CreateCategory)
				categoriesAdmin.PUT("/:id", categoryHandler.UpdateCategory)
				categoriesAdmin.DELETE("/:id", categoryHandler.DeleteCategory)
			}
		}

		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetAllProducts)
			products.GET("/featured", productHandler.GetFeaturedProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/slug/:slug", productHandler.GetProductBySlug)
			products.GET("/category/:category_id", productHandler.GetProductsByCategory)

			// Admin only
			productsAdmin := products.Group("")
			productsAdmin.Use(middleware.AuthMiddleware(cfg))
			productsAdmin.Use(middleware.AdminMiddleware())
			{
				productsAdmin.POST("", productHandler.CreateProduct)
				productsAdmin.PUT("/:id", productHandler.UpdateProduct)
				productsAdmin.DELETE("/:id", productHandler.DeleteProduct)
			}
		}

		// Cart routes (protected)
		cart := v1.Group("/cart")
		cart.Use(middleware.AuthMiddleware(cfg))
		{
			cart.POST("", cartHandler.AddToCart)
			cart.GET("", cartHandler.GetUserCart)
			cart.PUT("/:id", cartHandler.UpdateCartItem)
			cart.DELETE("/:id", cartHandler.RemoveFromCart)
			cart.DELETE("/clear", cartHandler.ClearCart)
		}

		// Wishlist routes (protected)
		wishlist := v1.Group("/wishlist")
		wishlist.Use(middleware.AuthMiddleware(cfg))
		{
			wishlist.POST("", wishlistHandler.AddToWishlist)
			wishlist.GET("", wishlistHandler.GetUserWishlist)
			wishlist.DELETE("/:id", wishlistHandler.RemoveFromWishlist)
			wishlist.DELETE("/clear", wishlistHandler.ClearWishlist)
		}

		// Order routes (protected)
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware(cfg))
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetUserOrders)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.POST("/:id/payment-proof", orderHandler.UploadPaymentProof)
			orders.POST("/:id/cancel", orderHandler.CancelOrder)
		}

		// Download routes (protected)
		downloads := v1.Group("/downloads")
		downloads.Use(middleware.AuthMiddleware(cfg))
		{
			downloads.POST("", downloadHandler.DownloadProduct)
			downloads.GET("", downloadHandler.GetUserDownloads)
			downloads.GET("/history/:product_id", downloadHandler.GetDownloadHistory)
		}

		// Review routes
		reviews := v1.Group("/reviews")
		{
			reviews.GET("/product/:product_id", reviewHandler.GetProductReviews)

			// Protected routes
			reviewsProtected := reviews.Group("")
			reviewsProtected.Use(middleware.AuthMiddleware(cfg))
			{
				reviewsProtected.POST("", reviewHandler.CreateReview)
				reviewsProtected.GET("/me", reviewHandler.GetMyReviews)
				reviewsProtected.PUT("/:id", reviewHandler.UpdateReview)
				reviewsProtected.DELETE("/:id", reviewHandler.DeleteReview)
			}
		}

		// Custom Order routes (protected)
		customOrders := v1.Group("/custom-orders")
		customOrders.Use(middleware.AuthMiddleware(cfg))
		{
			customOrders.POST("", customOrderHandler.CreateCustomOrder)
			customOrders.GET("/me", customOrderHandler.GetMyCustomOrders)
			customOrders.GET("/:id", customOrderHandler.GetCustomOrderByID)
			customOrders.PUT("/:id/cancel", customOrderHandler.CancelCustomOrder)
		}

		// Notification routes (protected)
		notifications := v1.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware(cfg))
		{
			notifications.GET("", notificationHandler.GetMyNotifications)
			notifications.GET("/unread", notificationHandler.GetUnreadNotifications)
			notifications.PUT("/:id/read", notificationHandler.MarkNotificationAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllNotificationsAsRead)
			notifications.DELETE("/:id", notificationHandler.DeleteNotification)
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

			// Orders management
			admin.GET("/orders", orderHandler.GetAllOrders)
			admin.POST("/orders/:id/approve", orderHandler.ApprovePayment)
			admin.POST("/orders/:id/reject", orderHandler.RejectPayment)

			// Custom Orders management
			admin.GET("/custom-orders", customOrderHandler.AdminGetAllCustomOrders)
			admin.PUT("/custom-orders/:id/process", customOrderHandler.AdminProcessCustomOrder)

			// Reviews moderation
			admin.DELETE("/reviews/:id", reviewHandler.AdminDeleteReview)

			// Analytics & Dashboard
			admin.GET("/analytics/dashboard", analyticsHandler.GetDashboardStats)
			admin.GET("/analytics/revenue", analyticsHandler.GetRevenueStats)
			admin.GET("/analytics/top-products", analyticsHandler.GetTopProducts)
			admin.GET("/analytics/users", analyticsHandler.GetUserStats)
			admin.GET("/analytics/orders", analyticsHandler.GetOrderStats)
		}
	}

	// Serve uploaded files
	r.Static("/uploads", "./uploads")
}
