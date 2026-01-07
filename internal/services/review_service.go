package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
)

type ReviewService interface {
	CreateReview(userID, productID uint, rating int, comment string) error
	GetProductReviews(productID uint, page, limit int) ([]*models.Review, int64, error)
	GetUserReviews(userID uint) ([]*models.Review, error)
	UpdateReview(reviewID, userID uint, rating int, comment string) error
	DeleteReview(reviewID, userID uint) error
	GetAverageRating(productID uint) (float64, error)
	AdminDeleteReview(reviewID uint) error
}

type reviewService struct {
	reviewRepo  repositories.ReviewRepository
	productRepo repositories.ProductRepository
	orderRepo   repositories.OrderRepository
}

func NewReviewService(
	reviewRepo repositories.ReviewRepository,
	productRepo repositories.ProductRepository,
	orderRepo repositories.OrderRepository,
) ReviewService {
	return &reviewService{
		reviewRepo:  reviewRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

func (s *reviewService) CreateReview(userID, productID uint, rating int, comment string) error {
	// Validate rating
	if rating < 1 || rating > 5 {
		return errors.New("rating harus antara 1-5")
	}

	// Check if product exists
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return errors.New("product tidak ditemukan")
	}

	// Check if user has purchased this product
	hasPurchased, err := s.orderRepo.HasUserPurchasedProduct(userID, productID)
	if err != nil {
		return err
	}
	if !hasPurchased {
		return errors.New("anda harus membeli product ini terlebih dahulu untuk memberikan review")
	}

	review := &models.Review{
		UserID:    userID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
	}

	return s.reviewRepo.Create(review)
}

func (s *reviewService) GetProductReviews(productID uint, page, limit int) ([]*models.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.reviewRepo.GetByProductID(productID, limit, offset)
}

func (s *reviewService) GetUserReviews(userID uint) ([]*models.Review, error) {
	return s.reviewRepo.GetByUserID(userID)
}

func (s *reviewService) UpdateReview(reviewID, userID uint, rating int, comment string) error {
	// Validate rating
	if rating < 1 || rating > 5 {
		return errors.New("rating harus antara 1-5")
	}

	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return errors.New("review tidak ditemukan")
	}

	// Check ownership
	if review.UserID != userID {
		return errors.New("anda tidak memiliki akses untuk mengubah review ini")
	}

	review.Rating = rating
	review.Comment = comment

	return s.reviewRepo.Update(review)
}

func (s *reviewService) DeleteReview(reviewID, userID uint) error {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return errors.New("review tidak ditemukan")
	}

	// Check ownership
	if review.UserID != userID {
		return errors.New("anda tidak memiliki akses untuk menghapus review ini")
	}

	return s.reviewRepo.Delete(reviewID)
}

func (s *reviewService) GetAverageRating(productID uint) (float64, error) {
	return s.reviewRepo.GetAverageRating(productID)
}

func (s *reviewService) AdminDeleteReview(reviewID uint) error {
	_, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return errors.New("review tidak ditemukan")
	}

	return s.reviewRepo.Delete(reviewID)
}
