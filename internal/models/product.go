package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"size:255;not null" json:"title"`
	Slug           string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Description    string         `gorm:"type:text" json:"description"`
	CategoryID     uint           `json:"category_id"`
	Category       *Category      `json:"category,omitempty"`
	Type           string         `gorm:"size:50;not null" json:"type"` // source_code, pdf, template, other
	Price          float64        `gorm:"not null" json:"price"`
	DiscountPrice  *float64       `json:"discount_price,omitempty"`
	PreviewImages  string         `gorm:"type:jsonb" json:"preview_images"` // JSON array
	DemoURL        string         `gorm:"size:500" json:"demo_url,omitempty"`
	FileURL        string         `gorm:"size:500" json:"file_url,omitempty"`
	TechStack      string         `gorm:"type:jsonb" json:"tech_stack"`   // JSON array
	Features       string         `gorm:"type:jsonb" json:"features"`     // JSON array
	Requirements   string         `gorm:"type:jsonb" json:"requirements"` // JSON array
	DownloadsCount int            `gorm:"default:0" json:"downloads_count"`
	ViewsCount     int            `gorm:"default:0" json:"views_count"`
	RatingAverage  float64        `gorm:"default:0" json:"rating_average"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedBy      uint           `json:"created_by"`
	Creator        *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Reviews        []Review       `json:"reviews,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductCreateRequest struct {
	Title         string   `json:"title" binding:"required"`
	Description   string   `json:"description" binding:"required"`
	CategoryID    uint     `json:"category_id" binding:"required"`
	Type          string   `json:"type" binding:"required,oneof=source_code pdf template other"`
	Price         float64  `json:"price" binding:"required,min=0"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
	DemoURL       string   `json:"demo_url,omitempty"`
	TechStack     []string `json:"tech_stack,omitempty"`
	Features      []string `json:"features,omitempty"`
	Requirements  []string `json:"requirements,omitempty"`
}

type ProductUpdateRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	CategoryID    uint     `json:"category_id"`
	Type          string   `json:"type" binding:"omitempty,oneof=source_code pdf template other"`
	Price         float64  `json:"price" binding:"min=0"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
	DemoURL       string   `json:"demo_url,omitempty"`
	TechStack     []string `json:"tech_stack,omitempty"`
	Features      []string `json:"features,omitempty"`
	Requirements  []string `json:"requirements,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

type ProductResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Description    string    `json:"description"`
	Category       *Category `json:"category,omitempty"`
	Type           string    `json:"type"`
	Price          float64   `json:"price"`
	DiscountPrice  *float64  `json:"discount_price,omitempty"`
	PreviewImages  []string  `json:"preview_images,omitempty"`
	DemoURL        string    `json:"demo_url,omitempty"`
	TechStack      []string  `json:"tech_stack,omitempty"`
	Features       []string  `json:"features,omitempty"`
	DownloadsCount int       `json:"downloads_count"`
	ViewsCount     int       `json:"views_count"`
	RatingAverage  float64   `json:"rating_average"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}
