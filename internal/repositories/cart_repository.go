package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *models.Cart) error
	GetByUserID(userID uint) ([]models.Cart, error)
	GetByID(id uint) (*models.Cart, error)
	GetByUserAndProduct(userID, productID uint) (*models.Cart, error)
	Update(cart *models.Cart) error
	Delete(id uint) error
	DeleteByUserID(userID uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) GetByUserID(userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Preload("Product").Preload("Product.Category").Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

func (r *cartRepository) GetByID(id uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Product").First(&cart, id).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetByUserAndProduct(userID, productID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Update(cart *models.Cart) error {
	return r.db.Save(cart).Error
}

func (r *cartRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cart{}, id).Error
}

func (r *cartRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Cart{}).Error
}
