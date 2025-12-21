package models

import (
	"time"
)

type Analytics struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Date        time.Time `gorm:"index;not null" json:"date"`
	MetricType  string    `gorm:"size:50;index;not null" json:"metric_type"` // page_view, unique_visitor, product_view, conversion, etc
	MetricValue int       `gorm:"not null" json:"metric_value"`
	Metadata    string    `gorm:"type:jsonb" json:"metadata"` // JSON - page, product_id, etc
	CreatedAt   time.Time `json:"created_at"`
}

type AnalyticsResponse struct {
	Date        time.Time              `json:"date"`
	MetricType  string                 `json:"metric_type"`
	MetricValue int                    `json:"metric_value"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type DashboardStats struct {
	TotalUsers      int     `json:"total_users"`
	ActiveUsers     int     `json:"active_users"`
	TotalProducts   int     `json:"total_products"`
	TotalOrders     int     `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	PendingOrders   int     `json:"pending_orders"`
	CompletedOrders int     `json:"completed_orders"`
	TotalDownloads  int     `json:"total_downloads"`
	TodayVisitors   int     `json:"today_visitors"`
	MonthlyRevenue  float64 `json:"monthly_revenue"`
}
