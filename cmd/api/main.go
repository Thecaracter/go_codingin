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

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, cfg)
	userHandler := handlers.NewUserHandler(userService)

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, cfg, authHandler, userHandler, apiLogRepo)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("📝 Environment: %s", cfg.AppEnv)
	log.Printf("🗄️  Database: Connected to %s", cfg.DBName)

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
