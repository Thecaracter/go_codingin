package models

import (
	"time"
)

type APILog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         *uint     `json:"user_id,omitempty"`
	Method         string    `gorm:"size:10;not null" json:"method"`
	Endpoint       string    `gorm:"size:500;not null" json:"endpoint"`
	StatusCode     int       `gorm:"not null" json:"status_code"`
	ResponseTimeMs int       `gorm:"not null" json:"response_time_ms"`
	IPAddress      string    `gorm:"size:45" json:"ip_address"`
	UserAgent      string    `gorm:"size:500" json:"user_agent"`
	RequestBody    string    `gorm:"type:text" json:"request_body,omitempty"`
	ResponseBody   string    `gorm:"type:text" json:"response_body,omitempty"`
	ErrorMessage   string    `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt      time.Time `gorm:"index" json:"created_at"`
}

type APILogResponse struct {
	ID             uint      `json:"id"`
	UserID         *uint     `json:"user_id,omitempty"`
	Method         string    `json:"method"`
	Endpoint       string    `json:"endpoint"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int       `json:"response_time_ms"`
	IPAddress      string    `json:"ip_address"`
	CreatedAt      time.Time `json:"created_at"`
}
