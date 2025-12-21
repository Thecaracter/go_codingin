package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	OrderID           uint           `json:"order_id"`
	Order             *Order         `json:"order,omitempty"`
	UserID            uint           `json:"user_id"`
	User              *User          `json:"user,omitempty"`
	TransactionNumber string         `gorm:"size:50;uniqueIndex;not null" json:"transaction_number"`
	Amount            float64        `gorm:"not null" json:"amount"`
	PaymentMethod     string         `gorm:"size:50;not null" json:"payment_method"`  // credit_card, bank_transfer, ewallet, crypto
	PaymentGateway    string         `gorm:"size:50;not null" json:"payment_gateway"` // midtrans, stripe, xendit
	PaymentGatewayRef string         `gorm:"size:255" json:"payment_gateway_ref"`
	Status            string         `gorm:"size:20;not null" json:"status"` // pending, success, failed, cancelled
	PaidAt            *time.Time     `json:"paid_at,omitempty"`
	Metadata          string         `gorm:"type:jsonb" json:"metadata"` // JSON
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type TransactionResponse struct {
	ID                uint       `json:"id"`
	TransactionNumber string     `json:"transaction_number"`
	OrderNumber       string     `json:"order_number"`
	Amount            float64    `json:"amount"`
	PaymentMethod     string     `json:"payment_method"`
	PaymentGateway    string     `json:"payment_gateway"`
	Status            string     `json:"status"`
	PaidAt            *time.Time `json:"paid_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}
