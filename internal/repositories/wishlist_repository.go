package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type WishlistRepository interface {
	Create(wishlist *models.Wishlist) error
	GetByUserID(userID uint) ([]models.Wishlist, error)
	GetByID(id uint) (*models.Wishlist, error)
	GetByUserAndProduct(userID, productID uint) (*models.Wishlist, error)
	Delete(id uint) error
	DeleteByUserID(userID uint) error
}

type wishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) Create(wishlist *models.Wishlist) error {
	return r.db.Create(wishlist).Error
}

func (r *wishlistRepository) GetByUserID(userID uint) ([]models.Wishlist, error) {
	var wishlists []models.Wishlist
	err := r.db.Preload("Product").Preload("Product.Category").Where("user_id = ?", userID).Find(&wishlists).Error
	return wishlists, err
}

func (r *wishlistRepository) GetByID(id uint) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.db.Preload("Product").First(&wishlist, id).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

func (r *wishlistRepository) GetByUserAndProduct(userID, productID uint) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&wishlist).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

func (r *wishlistRepository) Delete(id uint) error {
	return r.db.Delete(&models.Wishlist{}, id).Error
}

func (r *wishlistRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Wishlist{}).Error
}
