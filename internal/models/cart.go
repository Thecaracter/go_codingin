package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      *User          `json:"user,omitempty"`
	ProductID uint           `json:"product_id"`
	Product   *Product       `json:"product,omitempty"`
	Quantity  int            `gorm:"default:1" json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CartAddRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"min=1"`
}

type CartResponse struct {
	ID        uint            `json:"id"`
	Product   ProductResponse `json:"product"`
	Quantity  int             `json:"quantity"`
	CreatedAt time.Time       `json:"created_at"`
}
