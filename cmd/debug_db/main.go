package main

import (
	"fmt"
	"log"

	"gin-quickstart/internal/config"
	"gin-quickstart/internal/models"
	"gin-quickstart/pkg/utils"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("DEBUG: Connecting to DB_NAME=%s\n", cfg.DBName)
	db, err := utils.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🔍 Checking Products (Unscoped)...")
	var products []models.Product
	// Use Unscoped to see everything including deleted ones
	if err := db.Unscoped().Find(&products).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Printf("found %d products in total.\n", len(products))
	for _, p := range products {
		status := "Active"
		if p.DeletedAt.Valid {
			status = "DELETED"
		}
		fmt.Printf("- [%d] %s (%s) - %s\n", p.ID, p.Title, p.Slug, status)
	}
}
