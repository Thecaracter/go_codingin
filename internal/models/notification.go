package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      *User          `json:"user,omitempty"`
	Type      string         `gorm:"size:50;not null" json:"type"` // order, payment, system, custom_order
	Title     string         `gorm:"size:255;not null" json:"title"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	ActionURL string         `gorm:"size:500" json:"action_url,omitempty"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type NotificationResponse struct {
	ID        uint       `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	ActionURL string     `json:"action_url,omitempty"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
