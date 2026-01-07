package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type DownloadRepository interface {
	Create(download *models.Download) error
	GetByID(id uint) (*models.Download, error)
	GetByUserID(userID uint, page, limit int) ([]models.Download, int64, error)
	GetByUserAndProduct(userID, productID uint) ([]models.Download, error)
	CountByProduct(productID uint) (int64, error)
}

type downloadRepository struct {
	db *gorm.DB
}

func NewDownloadRepository(db *gorm.DB) DownloadRepository {
	return &downloadRepository{db: db}
}

func (r *downloadRepository) Create(download *models.Download) error {
	return r.db.Create(download).Error
}

func (r *downloadRepository) GetByID(id uint) (*models.Download, error) {
	var download models.Download
	err := r.db.Preload("User").Preload("Product").First(&download, id).Error
	if err != nil {
		return nil, err
	}
	return &download, nil
}

func (r *downloadRepository) GetByUserID(userID uint, page, limit int) ([]models.Download, int64, error) {
	var downloads []models.Download
	var total int64

	query := r.db.Model(&models.Download{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Product").Preload("Product.Category").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&downloads).Error

	return downloads, total, err
}

func (r *downloadRepository) GetByUserAndProduct(userID, productID uint) ([]models.Download, error) {
	var downloads []models.Download
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).
		Order("created_at DESC").
		Find(&downloads).Error
	return downloads, err
}

func (r *downloadRepository) CountByProduct(productID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Download{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}
