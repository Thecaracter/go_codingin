package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type APILogRepository interface {
	Create(log *models.APILog) error
	FindByDateRange(startDate, endDate string, limit int) ([]models.APILog, error)
	GetStatistics() (map[string]interface{}, error)
}

type apiLogRepository struct {
	db *gorm.DB
}

func NewAPILogRepository(db *gorm.DB) APILogRepository {
	return &apiLogRepository{db: db}
}

func (r *apiLogRepository) Create(log *models.APILog) error {
	return r.db.Create(log).Error
}

func (r *apiLogRepository) FindByDateRange(startDate, endDate string, limit int) ([]models.APILog, error) {
	var logs []models.APILog
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

func (r *apiLogRepository) GetStatistics() (map[string]interface{}, error) {
	var totalRequests int64
	var avgResponseTime float64

	r.db.Model(&models.APILog{}).Count(&totalRequests)
	r.db.Model(&models.APILog{}).Select("AVG(response_time_ms)").Scan(&avgResponseTime)

	return map[string]interface{}{
		"total_requests":    totalRequests,
		"avg_response_time": avgResponseTime,
	}, nil
}
