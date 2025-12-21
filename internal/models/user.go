package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Email      string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password   string         `gorm:"size:255" json:"-"` // Nullable for OAuth users
	Name       string         `gorm:"size:100;not null" json:"name"`
	Role       string         `gorm:"size:20;default:'user'" json:"role"`      // user, admin
	Provider   string         `gorm:"size:20;default:'local'" json:"provider"` // local, google, github
	ProviderID string         `gorm:"size:255" json:"provider_id,omitempty"`
	AvatarURL  string         `gorm:"size:500" json:"avatar_url,omitempty"`
	IsVerified bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserUpdateRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	Provider   string    `json:"provider"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
