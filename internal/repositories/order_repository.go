package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id uint) (*models.Order, error)
	GetByUserID(userID uint, page, limit int) ([]models.Order, int64, error)
	GetAll(page, limit int, status string) ([]models.Order, int64, error)
	Update(order *models.Order) error
	Delete(id uint) error
	GetByOrderNumber(orderNumber string) (*models.Order, error)
	HasUserPurchasedProduct(userID, productID uint) (bool, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("Product").Preload("Product.Category").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(userID uint, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Product").Preload("Product.Category").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) GetAll(page, limit int, status string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("User").Preload("Product").Preload("Product.Category").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}

func (r *orderRepository) GetByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("Product").Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) HasUserPurchasedProduct(userID, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Order{}).
		Where("user_id = ? AND product_id = ? AND order_status = ?", userID, productID, "completed").
		Count(&count).Error
	return count > 0, err
}
