package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetAll(page, limit int, categoryID *uint, search string) ([]models.Product, int64, error)
	GetByID(id uint) (*models.Product, error)
	GetBySlug(slug string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	GetFeatured(limit int) ([]models.Product, error)
	GetByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetAll(page, limit int, categoryID *uint, search string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Preload("Category")

	// Filter by category
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	// Search by name or description
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&products).Error

	return products, total, err
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySlug(slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) GetFeatured(limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("is_featured = ?", true).Limit(limit).Order("created_at DESC").Find(&products).Error
	return products, err
}

func (r *productRepository) GetByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Preload("Category").Where("category_id = ?", categoryID)

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&products).Error

	return products, total, err
}
