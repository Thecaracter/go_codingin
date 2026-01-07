package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type CustomOrderRepository interface {
	Create(customOrder *models.CustomOrder) error
	GetByID(id uint) (*models.CustomOrder, error)
	GetByUserID(userID uint) ([]*models.CustomOrder, error)
	GetAll(page, limit int, status string) ([]*models.CustomOrder, int64, error)
	Update(customOrder *models.CustomOrder) error
	Delete(id uint) error
}

type customOrderRepository struct {
	db *gorm.DB
}

func NewCustomOrderRepository(db *gorm.DB) CustomOrderRepository {
	return &customOrderRepository{db: db}
}

func (r *customOrderRepository) Create(customOrder *models.CustomOrder) error {
	return r.db.Create(customOrder).Error
}

func (r *customOrderRepository) GetByID(id uint) (*models.CustomOrder, error) {
	var customOrder models.CustomOrder
	err := r.db.Preload("User").First(&customOrder, id).Error
	if err != nil {
		return nil, err
	}
	return &customOrder, nil
}

func (r *customOrderRepository) GetByUserID(userID uint) ([]*models.CustomOrder, error) {
	var customOrders []*models.CustomOrder
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&customOrders).Error
	return customOrders, err
}

func (r *customOrderRepository) GetAll(page, limit int, status string) ([]*models.CustomOrder, int64, error) {
	var customOrders []*models.CustomOrder
	var total int64

	query := r.db.Model(&models.CustomOrder{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&customOrders).Error

	return customOrders, total, err
}

func (r *customOrderRepository) Update(customOrder *models.CustomOrder) error {
	return r.db.Save(customOrder).Error
}

func (r *customOrderRepository) Delete(id uint) error {
	return r.db.Delete(&models.CustomOrder{}, id).Error
}
