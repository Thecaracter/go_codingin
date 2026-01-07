package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"time"
)

type CustomOrderService interface {
	CreateRequest(userID uint, title, description, requirements string, budget float64) error
	GetUserRequests(userID uint) ([]*models.CustomOrder, error)
	GetRequestByID(id uint) (*models.CustomOrder, error)
	CancelRequest(id, userID uint) error

	// Admin
	GetAllRequests(page, limit int, status string) ([]*models.CustomOrder, int64, error)
	ProcessRequest(id uint, status, adminNotes string, quotedPrice *float64, estimatedDays *int) error
	CompleteRequest(id uint) error
}

type customOrderService struct {
	customOrderRepo repositories.CustomOrderRepository
}

func NewCustomOrderService(customOrderRepo repositories.CustomOrderRepository) CustomOrderService {
	return &customOrderService{
		customOrderRepo: customOrderRepo,
	}
}

func (s *customOrderService) CreateRequest(userID uint, title, description, requirements string, budget float64) error {
	if title == "" || description == "" {
		return errors.New("title dan description harus diisi")
	}

	if budget < 0 {
		return errors.New("budget tidak valid")
	}

	customOrder := &models.CustomOrder{
		UserID:       userID,
		Title:        title,
		Description:  description,
		Requirements: requirements,
		BudgetMin:    budget,
		BudgetMax:    budget,
		Status:       "pending",
	}

	return s.customOrderRepo.Create(customOrder)
}

func (s *customOrderService) GetUserRequests(userID uint) ([]*models.CustomOrder, error) {
	return s.customOrderRepo.GetByUserID(userID)
}

func (s *customOrderService) GetRequestByID(id uint) (*models.CustomOrder, error) {
	return s.customOrderRepo.GetByID(id)
}

func (s *customOrderService) CancelRequest(id, userID uint) error {
	customOrder, err := s.customOrderRepo.GetByID(id)
	if err != nil {
		return errors.New("custom order tidak ditemukan")
	}

	if customOrder.UserID != userID {
		return errors.New("anda tidak memiliki akses untuk membatalkan request ini")
	}

	if customOrder.Status == "completed" || customOrder.Status == "cancelled" {
		return errors.New("custom order ini tidak dapat dibatalkan")
	}

	customOrder.Status = "cancelled"
	return s.customOrderRepo.Update(customOrder)
}

func (s *customOrderService) GetAllRequests(page, limit int, status string) ([]*models.CustomOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.customOrderRepo.GetAll(page, limit, status)
}

func (s *customOrderService) ProcessRequest(id uint, status, adminNotes string, quotedPrice *float64, estimatedDays *int) error {
	customOrder, err := s.customOrderRepo.GetByID(id)
	if err != nil {
		return errors.New("custom order tidak ditemukan")
	}

	if customOrder.Status == "completed" || customOrder.Status == "cancelled" {
		return errors.New("custom order ini sudah selesai atau dibatalkan")
	}

	validStatuses := map[string]bool{
		"pending":     true,
		"reviewing":   true,
		"quoted":      true,
		"in_progress": true,
		"completed":   true,
		"cancelled":   true,
	}

	if !validStatuses[status] {
		return errors.New("status tidak valid")
	}

	customOrder.Status = status
	customOrder.AdminNotes = adminNotes

	if quotedPrice != nil {
		if *quotedPrice < 0 {
			return errors.New("harga tidak valid")
		}
		customOrder.QuoteAmount = quotedPrice
	}

	if estimatedDays != nil {
		if *estimatedDays < 0 {
			return errors.New("estimasi hari tidak valid")
		}
		deadline := time.Now().AddDate(0, 0, *estimatedDays)
		customOrder.Deadline = &deadline
	}

	return s.customOrderRepo.Update(customOrder)
}

func (s *customOrderService) CompleteRequest(id uint) error {
	customOrder, err := s.customOrderRepo.GetByID(id)
	if err != nil {
		return errors.New("custom order tidak ditemukan")
	}

	if customOrder.Status == "completed" {
		return errors.New("custom order ini sudah selesai")
	}

	if customOrder.Status == "cancelled" {
		return errors.New("custom order ini sudah dibatalkan")
	}

	customOrder.Status = "completed"

	return s.customOrderRepo.Update(customOrder)
}
