package models

import (
	"time"

	"gorm.io/gorm"
)

type Wishlist struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      *User          `json:"user,omitempty"`
	ProductID uint           `json:"product_id"`
	Product   *Product       `json:"product,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type WishlistAddRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

type WishlistResponse struct {
	ID        uint            `json:"id"`
	Product   ProductResponse `json:"product"`
	CreatedAt time.Time       `json:"created_at"`
}
