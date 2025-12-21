package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Slug        string         `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"size:255" json:"icon,omitempty"`
	ParentID    *uint          `json:"parent_id,omitempty"`
	Parent      *Category      `json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Order       int            `gorm:"default:0" json:"order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Products    []Product      `json:"products,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ParentID    *uint  `json:"parent_id"`
	Order       int    `json:"order"`
	IsActive    *bool  `json:"is_active"`
}

type CategoryResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Description string             `json:"description"`
	Icon        string             `json:"icon,omitempty"`
	ParentID    *uint              `json:"parent_id,omitempty"`
	Children    []CategoryResponse `json:"children,omitempty"`
	Order       int                `json:"order"`
	IsActive    bool               `json:"is_active"`
	CreatedAt   time.Time          `json:"created_at"`
}
