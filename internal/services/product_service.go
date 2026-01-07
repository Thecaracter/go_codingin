package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"gin-quickstart/pkg/utils"
	"mime/multipart"
	"strings"
	"unicode"
)

type ProductService interface {
	CreateProduct(req CreateProductRequest, file *multipart.FileHeader, createdBy uint) (*models.Product, error)
	GetAllProducts(page, limit int, categoryID *uint, search string) ([]models.Product, int64, error)
	GetProductByID(id uint) (*models.Product, error)
	GetProductBySlug(slug string) (*models.Product, error)
	UpdateProduct(id uint, req UpdateProductRequest, file *multipart.FileHeader) (*models.Product, error)
	DeleteProduct(id uint) error
	GetFeaturedProducts(limit int) ([]models.Product, error)
	GetProductsByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error)
}

type CreateProductRequest struct {
	Title         string   `json:"title" binding:"required"`
	Description   string   `json:"description"`
	Type          string   `json:"type" binding:"required"`
	Price         float64  `json:"price" binding:"required,gt=0"`
	DiscountPrice *float64 `json:"discount_price"`
	CategoryID    uint     `json:"category_id" binding:"required"`
	DemoURL       string   `json:"demo_url"`
	TechStack     []string `json:"tech_stack"`
	Features      []string `json:"features"`
	Requirements  []string `json:"requirements"`
	IsActive      bool     `json:"is_active"`
}

type UpdateProductRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Type          string   `json:"type"`
	Price         float64  `json:"price"`
	DiscountPrice *float64 `json:"discount_price"`
	CategoryID    uint     `json:"category_id"`
	DemoURL       string   `json:"demo_url"`
	TechStack     []string `json:"tech_stack"`
	Features      []string `json:"features"`
	Requirements  []string `json:"requirements"`
	IsActive      bool     `json:"is_active"`
}

type productService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(productRepo repositories.ProductRepository, categoryRepo repositories.CategoryRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) CreateProduct(req CreateProductRequest, file *multipart.FileHeader, createdBy uint) (*models.Product, error) {
	_, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	slug := generateProductSlug(req.Title)

	existing, _ := s.productRepo.GetBySlug(slug)
	if existing != nil {
		return nil, errors.New("product with this title already exists")
	}

	product := &models.Product{
		Title:       req.Title,
		Slug:        slug,
		Description: req.Description,
		Type:        req.Type,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		IsActive:    req.IsActive,
		CreatedBy:   createdBy,
	}

	if req.DiscountPrice != nil {
		product.DiscountPrice = req.DiscountPrice
	}
	if req.DemoURL != "" {
		product.DemoURL = req.DemoURL
	}

	var imagePath string
	if file != nil {
		var err error
		imagePath, err = utils.UploadFile(file, "products")
		if err != nil {
			return nil, err
		}
		product.PreviewImages = `["` + imagePath + `"]`
	}

	if err := s.productRepo.Create(product); err != nil {
		if imagePath != "" {
			utils.DeleteFile(imagePath)
		}
		return nil, err
	}

	return product, nil
}

func (s *productService) GetAllProducts(page, limit int, categoryID *uint, search string) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return s.productRepo.GetAll(page, limit, categoryID, search)
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *productService) GetProductBySlug(slug string) (*models.Product, error) {
	product, err := s.productRepo.GetBySlug(slug)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *productService) UpdateProduct(id uint, req UpdateProductRequest, file *multipart.FileHeader) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if req.Title != "" {
		product.Title = req.Title
		product.Slug = generateProductSlug(req.Title)
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Type != "" {
		product.Type = req.Type
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.DiscountPrice != nil {
		product.DiscountPrice = req.DiscountPrice
	}
	if req.CategoryID > 0 {
		_, err := s.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = req.CategoryID
	}
	if req.DemoURL != "" {
		product.DemoURL = req.DemoURL
	}

	product.IsActive = req.IsActive

	if file != nil {
		if product.PreviewImages != "" {
			// Simple deletion - in real app would parse JSON properly
		}

		imagePath, err := utils.UploadFile(file, "products")
		if err != nil {
			return nil, err
		}
		product.PreviewImages = `["` + imagePath + `"]`
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(id uint) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	if product.PreviewImages != "" {
	}

	return s.productRepo.Delete(id)
}

func (s *productService) GetFeaturedProducts(limit int) ([]models.Product, error) {
	if limit < 1 {
		limit = 10
	}
	products, _, err := s.productRepo.GetAll(1, limit, nil, "")
	return products, err
}

func (s *productService) GetProductsByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Validate category exists
	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, 0, errors.New("category not found")
	}

	return s.productRepo.GetByCategory(categoryID, page, limit)
}

// generateProductSlug creates URL-friendly slug from name
func generateProductSlug(name string) string {
	slug := strings.ToLower(name)

	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) {
			result.WriteRune('-')
		}
	}

	slug = result.String()
	slug = strings.ReplaceAll(slug, "--", "-")
	slug = strings.Trim(slug, "-")

	return slug
}
