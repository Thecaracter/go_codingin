package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `json:"user_id"`
	User               *User          `json:"user,omitempty"`
	ProductID          uint           `json:"product_id"`
	Product            *Product       `json:"product,omitempty"`
	OrderID            uint           `json:"order_id"`
	Order              *Order         `json:"order,omitempty"`
	Rating             int            `gorm:"not null" json:"rating"` // 1-5
	Comment            string         `gorm:"type:text" json:"comment"`
	IsVerifiedPurchase bool           `gorm:"default:false" json:"is_verified_purchase"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

type ReviewCreateRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	OrderID   uint   `json:"order_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required"`
}

type ReviewUpdateRequest struct {
	Rating  int    `json:"rating" binding:"min=1,max=5"`
	Comment string `json:"comment"`
}

type ReviewResponse struct {
	ID                 uint          `json:"id"`
	User               *UserResponse `json:"user"`
	Rating             int           `json:"rating"`
	Comment            string        `json:"comment"`
	IsVerifiedPurchase bool          `json:"is_verified_purchase"`
	CreatedAt          time.Time     `json:"created_at"`
}
