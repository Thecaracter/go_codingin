package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomOrder struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderID      uint           `gorm:"uniqueIndex" json:"order_id"`
	Order        *Order         `json:"order,omitempty"`
	UserID       uint           `json:"user_id"`
	User         *User          `json:"user,omitempty"`
	Title        string         `gorm:"size:255;not null" json:"title"`
	Description  string         `gorm:"type:text;not null" json:"description"`
	Requirements string         `gorm:"type:text" json:"requirements"`
	BudgetMin    float64        `json:"budget_min"`
	BudgetMax    float64        `json:"budget_max"`
	Deadline     *time.Time     `json:"deadline,omitempty"`
	Status       string         `gorm:"size:30;not null" json:"status"` // submitted, under_review, in_progress, completed, cancelled
	AdminNotes   string         `gorm:"type:text" json:"admin_notes,omitempty"`
	Attachments  string         `gorm:"type:jsonb" json:"attachments"` // JSON array
	QuoteAmount  *float64       `json:"quote_amount,omitempty"`
	AgreedAmount *float64       `json:"agreed_amount,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type CustomOrderCreateRequest struct {
	Title        string     `json:"title" binding:"required"`
	Description  string     `json:"description" binding:"required"`
	Requirements string     `json:"requirements"`
	BudgetMin    float64    `json:"budget_min" binding:"required,min=0"`
	BudgetMax    float64    `json:"budget_max" binding:"required,min=0"`
	Deadline     *time.Time `json:"deadline,omitempty"`
}

type CustomOrderUpdateRequest struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Requirements string     `json:"requirements"`
	BudgetMin    float64    `json:"budget_min" binding:"min=0"`
	BudgetMax    float64    `json:"budget_max" binding:"min=0"`
	Deadline     *time.Time `json:"deadline,omitempty"`
}

type CustomOrderAdminUpdateRequest struct {
	Status       string   `json:"status" binding:"omitempty,oneof=submitted under_review in_progress completed cancelled"`
	AdminNotes   string   `json:"admin_notes"`
	QuoteAmount  *float64 `json:"quote_amount"`
	AgreedAmount *float64 `json:"agreed_amount"`
}

type CustomOrderResponse struct {
	ID           uint       `json:"id"`
	OrderNumber  string     `json:"order_number"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Requirements string     `json:"requirements"`
	BudgetMin    float64    `json:"budget_min"`
	BudgetMax    float64    `json:"budget_max"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	Status       string     `json:"status"`
	AdminNotes   string     `json:"admin_notes,omitempty"`
	QuoteAmount  *float64   `json:"quote_amount,omitempty"`
	AgreedAmount *float64   `json:"agreed_amount,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
