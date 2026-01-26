package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gin-quickstart/internal/config"
	"gin-quickstart/internal/models"
	"gin-quickstart/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	wd, _ := os.Getwd()
	fmt.Printf("DEBUG_CONFIG: DB_NAME=%s CWD=%s\n", cfg.DBName, wd)

	// Initialize database
	db, err := utils.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("🌱 Seeding database...")

	// 1. Create Users
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	admin := models.User{
		Name:       "Admin User",
		Email:      "admin@codingin.com",
		Password:   string(hashedPassword),
		Role:       "admin",
		IsVerified: true,
		Provider:   "local",
	}

	if err := db.Where(models.User{Email: admin.Email}).FirstOrCreate(&admin).Error; err != nil {
		log.Fatal("Failed to seed admin:", err)
	}
	fmt.Printf("✅ User seeded: %s\n", admin.Email)

	// 2. Create Categories
	categories := []models.Category{
		{
			Name:        "Source Code",
			Slug:        "source-code",
			Description: "Full application source codes",
			Icon:        "code",
			Order:       1,
		},
		{
			Name:        "Templates",
			Slug:        "templates",
			Description: "UI kits and website templates",
			Icon:        "layout",
			Order:       2,
		},
		{
			Name:        "Mobile Apps",
			Slug:        "mobile-apps",
			Description: "Flutter, React Native, and native apps",
			Icon:        "smartphone",
			Order:       3,
		},
		{
			Name:        "Scripts",
			Slug:        "scripts",
			Description: "Python, PHP, and JS utility scripts",
			Icon:        "terminal",
			Order:       4,
		},
	}

	var createdCategories []models.Category
	for _, cat := range categories {
		if err := db.Where(models.Category{Slug: cat.Slug}).FirstOrCreate(&cat).Error; err != nil {
			log.Printf("Failed to seed category %s: %v", cat.Name, err)
		} else {
			createdCategories = append(createdCategories, cat)
			fmt.Printf("✅ Category seeded: %s\n", cat.Name)
		}
	}

	// 3. Create Products dengan JSON yang benar
	products := []struct {
		Title         string
		Slug          string
		Description   string
		Type          string
		Price         float64
		PreviewImages []string
		TechStack     []string
		Features      []string
		CategoryID    uint
	}{
		{
			Title:         "E-Commerce Microservices Go",
			Slug:          "ecommerce-microservices-go",
			Description:   "A complete microservices-based e-commerce backend written in Go (Golang). Includes services for products, orders, users, and payments. Features gRPC communication, Kafka event bus, and Docker compose setup.",
			Type:          "source_code",
			Price:         1500000,
			PreviewImages: []string{"https://images.unsplash.com/photo-1555949963-aa79dcee981c?auto=format&fit=crop&q=80&w=1000", "https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&q=80&w=1000"},
			TechStack:     []string{"Golang", "gRPC", "Kafka", "PostgreSQL", "Docker"},
			Features:      []string{"Microservices Architecture", "gRPC Communication", "Event-Driven", "Scalable"},
			CategoryID:    createdCategories[0].ID,
		},
		{
			Title:         "NextJS SaaS Boilerplate",
			Slug:          "nextjs-saas-boilerplate",
			Description:   "Launch your SaaS in days with this comprehensive Next.js boilerplate. Includes Stripe integration, Authentication, Database setup with Prisma, and tailwind UI components.",
			Type:          "source_code",
			Price:         750000,
			PreviewImages: []string{"https://images.unsplash.com/photo-1618477247222-ac591245363d?auto=format&fit=crop&q=80&w=1000"},
			TechStack:     []string{"Next.js", "TypeScript", "Tailwind CSS", "Prisma", "Stripe"},
			Features:      []string{"Authentication", "Subscription Payments", "Admin Dashboard", "SEO Optimized"},
			CategoryID:    createdCategories[0].ID,
		},
		{
			Title:         "Flutter Delivery App",
			Slug:          "flutter-delivery-app",
			Description:   "Uber Eats clone made with Flutter. Beautiful animations, real-time tracking, and vendor dashboard included. Works on both iOS and Android.",
			Type:          "mobile_apps",
			Price:         500000,
			PreviewImages: []string{"https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?auto=format&fit=crop&q=80&w=1000"},
			TechStack:     []string{"Flutter", "Firebase", "Google Maps API"},
			Features:      []string{"Real-time Tracking", "Push Notifications", "Vendor App", "Dark Mode"},
			CategoryID:    createdCategories[2].ID,
		},
		{
			Title:         "Modern Admin Dashboard UI",
			Slug:          "modern-admin-dashboard",
			Description:   "Clean, responsive, and customizable admin dashboard template built with Vue 3 and Tailwind CSS.",
			Type:          "template",
			Price:         250000,
			PreviewImages: []string{"https://images.unsplash.com/photo-1551288049-bebda4e38f71?auto=format&fit=crop&q=80&w=1000"},
			TechStack:     []string{"Vue 3", "Tailwind CSS", "Vite"},
			Features:      []string{"Dark Mode Support", "Responsive Layout", "Chart.js Integration", "Form Components"},
			CategoryID:    createdCategories[1].ID,
		},
		{
			Title:         "Instagram Automation Bot",
			Slug:          "instagram-automation-bot",
			Description:   "Python script to automate likes, follows, and comments. Includes safety limits to prevent account bans.",
			Type:          "scripts",
			Price:         150000,
			PreviewImages: []string{"https://images.unsplash.com/photo-1611162617474-5b21e879e113?auto=format&fit=crop&q=80&w=1000"},
			TechStack:     []string{"Python", "Selenium"},
			Features:      []string{"Auto Like", "Auto Follow", "Targeting by Hashtag"},
			CategoryID:    createdCategories[3].ID,
		},
	}

	for _, seedP := range products {
		// Convert arrays to JSON strings
		previewImagesJSON, _ := json.Marshal(seedP.PreviewImages)
		techStackJSON, _ := json.Marshal(seedP.TechStack)
		featuresJSON, _ := json.Marshal(seedP.Features)

		product := models.Product{
			Title:         seedP.Title,
			Slug:          seedP.Slug,
			Description:   seedP.Description,
			CategoryID:    seedP.CategoryID,
			Type:          seedP.Type,
			Price:         seedP.Price,
			PreviewImages: string(previewImagesJSON),
			TechStack:     string(techStackJSON),
			Features:      string(featuresJSON),
			Requirements:  "[]", // Empty JSON array
			IsActive:      true,
			CreatedBy:     admin.ID,
		}

		// Check if product exists
		var existingProduct models.Product
		err := db.Where("slug = ?", product.Slug).First(&existingProduct).Error

		if err != nil {
			// Product doesn't exist, create it
			if err := db.Create(&product).Error; err != nil {
				log.Printf("❌ Failed to seed product %s: %v", product.Title, err)
			} else {
				fmt.Printf("✅ Product seeded: %s\n", product.Title)
			}
		} else {
			fmt.Printf("⏭️  Product already exists: %s\n", product.Title)
		}
	}

	fmt.Println("🚀 Database seeding completed!")
}
