package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrderNumber    string         `gorm:"size:50;uniqueIndex;not null" json:"order_number"`
	UserID         uint           `json:"user_id"`
	User           *User          `json:"user,omitempty"`
	ProductID      *uint          `json:"product_id,omitempty"`
	Product        *Product       `json:"product,omitempty"`
	OrderType      string         `gorm:"size:20;not null" json:"order_type"` // product, custom
	Status         string         `gorm:"size:20;not null" json:"status"`     // pending, processing, completed, cancelled, refunded
	TotalAmount    float64        `gorm:"not null" json:"total_amount"`
	DiscountAmount float64        `gorm:"default:0" json:"discount_amount"`
	FinalAmount    float64        `gorm:"not null" json:"final_amount"`
	PaymentMethod  string         `gorm:"size:50" json:"payment_method"`
	PaymentStatus  string         `gorm:"size:20;not null" json:"payment_status"` // pending, paid, failed, refunded
	PaymentID      string         `gorm:"size:255" json:"payment_id"`
	Notes          string         `gorm:"type:text" json:"notes,omitempty"`
	CustomOrder    *CustomOrder   `json:"custom_order,omitempty"`
	Transactions   []Transaction  `json:"transactions,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderCreateRequest struct {
	ProductID     *uint  `json:"product_id"`
	OrderType     string `json:"order_type" binding:"required,oneof=product custom"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Notes         string `json:"notes,omitempty"`
}

type OrderResponse struct {
	ID             uint         `json:"id"`
	OrderNumber    string       `json:"order_number"`
	OrderType      string       `json:"order_type"`
	Status         string       `json:"status"`
	Product        *Product     `json:"product,omitempty"`
	TotalAmount    float64      `json:"total_amount"`
	DiscountAmount float64      `json:"discount_amount"`
	FinalAmount    float64      `json:"final_amount"`
	PaymentMethod  string       `json:"payment_method"`
	PaymentStatus  string       `json:"payment_status"`
	CustomOrder    *CustomOrder `json:"custom_order,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}
