package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"strings"
	"unicode"
)

type CategoryService interface {
	CreateCategory(name, description string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	GetCategoryByID(id uint) (*models.Category, error)
	GetCategoryBySlug(slug string) (*models.Category, error)
	UpdateCategory(id uint, name, description string) (*models.Category, error)
	DeleteCategory(id uint) error
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(name, description string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}

	slug := generateSlug(name)

	// Check if slug already exists
	existing, _ := s.categoryRepo.GetBySlug(slug)
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	category := &models.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) GetCategoryByID(id uint) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryService) GetCategoryBySlug(slug string) (*models.Category, error) {
	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryService) UpdateCategory(id uint, name, description string) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if name != "" {
		category.Name = name
		category.Slug = generateSlug(name)
	}

	if description != "" {
		category.Description = description
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	return s.categoryRepo.Delete(category.ID)
}

// generateSlug creates URL-friendly slug from name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and special chars with dash
	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) {
			result.WriteRune('-')
		}
	}

	// Remove multiple consecutive dashes
	slug = result.String()
	slug = strings.ReplaceAll(slug, "--", "-")
	slug = strings.Trim(slug, "-")

	return slug
}
