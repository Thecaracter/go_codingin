package main

import (
	"fmt"
	"log"

	"gin-quickstart/internal/config"
	"gin-quickstart/internal/handlers"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"gin-quickstart/internal/routes"
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := utils.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate all models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
		&models.Order{},
		&models.CustomOrder{},
		&models.Transaction{},
		&models.Download{},
		&models.Review{},
		&models.Cart{},
		&models.Wishlist{},
		&models.APILog{},
		&models.Analytics{},
		&models.Notification{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully")

	// Set Gin mode based on environment
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	apiLogRepo := repositories.NewAPILogRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	wishlistRepo := repositories.NewWishlistRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	downloadRepo := repositories.NewDownloadRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	customOrderRepo := repositories.NewCustomOrderRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	wishlistService := services.NewWishlistService(wishlistRepo, productRepo)
	orderService := services.NewOrderService(orderRepo, transactionRepo, productRepo, cartRepo)
	downloadService := services.NewDownloadService(downloadRepo, orderRepo, productRepo)
	reviewService := services.NewReviewService(reviewRepo, productRepo, orderRepo)
	customOrderService := services.NewCustomOrderService(customOrderRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	analyticsService := services.NewAnalyticsService(userRepo, productRepo, orderRepo, reviewRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, cfg)
	userHandler := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
	orderHandler := handlers.NewOrderHandler(orderService)
	downloadHandler := handlers.NewDownloadHandler(downloadService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	customOrderHandler := handlers.NewCustomOrderHandler(customOrderService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, cfg, authHandler, userHandler, categoryHandler, productHandler, cartHandler, wishlistHandler, orderHandler, downloadHandler, reviewHandler, customOrderHandler, notificationHandler, analyticsHandler, apiLogRepo)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("📝 Environment: %s", cfg.AppEnv)
	log.Printf("🗄️  Database: Connected to %s", cfg.DBName)

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
