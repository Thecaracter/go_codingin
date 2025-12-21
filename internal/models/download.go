package models

import (
	"time"

	"gorm.io/gorm"
)

type Download struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `json:"user_id"`
	User         *User          `json:"user,omitempty"`
	ProductID    uint           `json:"product_id"`
	Product      *Product       `json:"product,omitempty"`
	OrderID      uint           `json:"order_id"`
	Order        *Order         `json:"order,omitempty"`
	DownloadURL  string         `gorm:"size:500" json:"download_url"` // Temporary signed URL
	ExpiresAt    time.Time      `json:"expires_at"`
	IsUsed       bool           `gorm:"default:false" json:"is_used"`
	DownloadedAt *time.Time     `json:"downloaded_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type DownloadResponse struct {
	ID          uint      `json:"id"`
	Product     *Product  `json:"product"`
	DownloadURL string    `json:"download_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}
