package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
)

type DownloadService interface {
	CreateDownload(userID, orderID, productID uint) error
	GetUserDownloads(userID uint, page, limit int) ([]models.Download, int64, error)
	CanDownload(userID, productID uint) (bool, error)
	GetDownloadHistory(userID, productID uint) ([]models.Download, error)
}

type downloadService struct {
	downloadRepo repositories.DownloadRepository
	orderRepo    repositories.OrderRepository
	productRepo  repositories.ProductRepository
}

func NewDownloadService(
	downloadRepo repositories.DownloadRepository,
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
) DownloadService {
	return &downloadService{
		downloadRepo: downloadRepo,
		orderRepo:    orderRepo,
		productRepo:  productRepo,
	}
}

func (s *downloadService) CreateDownload(userID, orderID, productID uint) error {
	// Verify order exists and belongs to user
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("unauthorized")
	}

	if order.Status != "paid" {
		return errors.New("order is not paid")
	}

	if order.ProductID == nil || *order.ProductID != productID {
		return errors.New("product does not match order")
	}

	// Create download record
	download := &models.Download{
		UserID:    userID,
		ProductID: productID,
		OrderID:   orderID,
	}

	if err := s.downloadRepo.Create(download); err != nil {
		return err
	}

	// Update product download count
	product, err := s.productRepo.GetByID(productID)
	if err == nil {
		product.DownloadsCount++
		s.productRepo.Update(product)
	}

	return nil
}

func (s *downloadService) GetUserDownloads(userID uint, page, limit int) ([]models.Download, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return s.downloadRepo.GetByUserID(userID, page, limit)
}

func (s *downloadService) CanDownload(userID, productID uint) (bool, error) {
	// Check if user has paid order for this product
	orders, _, err := s.orderRepo.GetByUserID(userID, 1, 100)
	if err != nil {
		return false, err
	}

	for _, order := range orders {
		if order.ProductID != nil && *order.ProductID == productID && order.PaymentStatus == "paid" {
			return true, nil
		}
	}

	return false, nil
}

func (s *downloadService) GetDownloadHistory(userID, productID uint) ([]models.Download, error) {
	return s.downloadRepo.GetByUserAndProduct(userID, productID)
}
