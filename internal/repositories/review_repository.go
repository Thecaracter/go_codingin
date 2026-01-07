package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	GetByID(id uint) (*models.Review, error)
	GetByProductID(productID uint, limit, offset int) ([]*models.Review, int64, error)
	GetByUserID(userID uint) ([]*models.Review, error)
	Update(review *models.Review) error
	Delete(id uint) error
	GetAverageRating(productID uint) (float64, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User").Preload("Product").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByProductID(productID uint, limit, offset int) ([]*models.Review, int64, error) {
	var reviews []*models.Review
	var total int64

	query := r.db.Model(&models.Review{}).Where("product_id = ?", productID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error

	return reviews, total, err
}

func (r *reviewRepository) GetByUserID(userID uint) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.Where("user_id = ?", userID).
		Preload("Product").
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

func (r *reviewRepository) GetAverageRating(productID uint) (float64, error) {
	var avg float64
	err := r.db.Model(&models.Review{}).
		Where("product_id = ?", productID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avg).Error
	return avg, err
}
