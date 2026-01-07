package services

import (
	"errors"
	"fmt"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"gin-quickstart/pkg/utils"
	"mime/multipart"
	"time"
)

type OrderService interface {
	CreateOrder(userID, productID uint, quantity int) (*models.Order, error)
	GetOrderByID(userID, orderID uint) (*models.Order, error)
	GetUserOrders(userID uint, page, limit int) ([]models.Order, int64, error)
	GetAllOrders(page, limit int, status string) ([]models.Order, int64, error)
	UploadPaymentProof(userID, orderID uint, file *multipart.FileHeader) error
	ApprovePayment(orderID uint, adminID uint) error
	RejectPayment(orderID uint, adminID uint, reason string) error
	CancelOrder(userID, orderID uint) error
}

type orderService struct {
	orderRepo       repositories.OrderRepository
	transactionRepo repositories.TransactionRepository
	productRepo     repositories.ProductRepository
	cartRepo        repositories.CartRepository
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	transactionRepo repositories.TransactionRepository,
	productRepo repositories.ProductRepository,
	cartRepo repositories.CartRepository,
) OrderService {
	return &orderService{
		orderRepo:       orderRepo,
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
		cartRepo:        cartRepo,
	}
}

func (s *orderService) CreateOrder(userID, productID uint, quantity int) (*models.Order, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if !product.IsActive {
		return nil, errors.New("product is not available")
	}

	var totalAmount float64
	if product.DiscountPrice != nil && *product.DiscountPrice > 0 {
		totalAmount = *product.DiscountPrice * float64(quantity)
	} else {
		totalAmount = product.Price * float64(quantity)
	}

	orderNumber := fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), userID)

	order := &models.Order{
		OrderNumber:   orderNumber,
		UserID:        userID,
		ProductID:     &productID,
		OrderType:     "product",
		TotalAmount:   totalAmount,
		FinalAmount:   totalAmount,
		Status:        "pending",
		PaymentMethod: "manual_transfer",
		PaymentStatus: "pending",
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	transaction := &models.Transaction{
		OrderID:       order.ID,
		Amount:        totalAmount,
		Status:        "pending",
		PaymentMethod: "manual_transfer",
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrderByID(userID, orderID uint) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return order, nil
}

func (s *orderService) GetUserOrders(userID uint, page, limit int) ([]models.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return s.orderRepo.GetByUserID(userID, page, limit)
}

func (s *orderService) GetAllOrders(page, limit int, status string) ([]models.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return s.orderRepo.GetAll(page, limit, status)
}

func (s *orderService) UploadPaymentProof(userID, orderID uint, file *multipart.FileHeader) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("unauthorized")
	}

	if order.PaymentStatus != "pending" {
		return errors.New("order is not pending payment")
	}

	filePath, err := utils.UploadFile(file, "payment_proofs")
	if err != nil {
		return err
	}

	// Get transaction
	transaction, err := s.transactionRepo.GetByOrderID(order.ID)
	if err != nil {
		utils.DeleteFile(filePath)
		return errors.New("transaction not found")
	}

	transaction.Metadata = `{"proof_image":"` + filePath + `"}`
	transaction.Status = "pending"

	if err := s.transactionRepo.Update(transaction); err != nil {
		utils.DeleteFile(filePath)
		return err
	}

	order.Status = "processing"
	order.PaymentStatus = "pending"
	if err := s.orderRepo.Update(order); err != nil {
		return err
	}

	return nil
}

func (s *orderService) ApprovePayment(orderID uint, adminID uint) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.PaymentStatus != "pending" || order.Status != "processing" {
		return errors.New("order is not waiting for confirmation")
	}

	// Get transaction
	transaction, err := s.transactionRepo.GetByOrderID(order.ID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Update transaction
	now := time.Now()
	transaction.Status = "success"
	transaction.PaidAt = &now
	// Store admin info in Metadata
	transaction.Metadata = fmt.Sprintf(`{"verified_by":%d,"verified_at":"%s"}`, adminID, now.Format(time.RFC3339))

	if err := s.transactionRepo.Update(transaction); err != nil {
		return err
	}

	// Update order
	order.Status = "completed"
	order.PaymentStatus = "paid"
	if err := s.orderRepo.Update(order); err != nil {
		return err
	}

	return nil
}

func (s *orderService) RejectPayment(orderID uint, adminID uint, reason string) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.PaymentStatus != "pending" || order.Status != "processing" {
		return errors.New("order is not waiting for confirmation")
	}

	// Get transaction
	transaction, err := s.transactionRepo.GetByOrderID(order.ID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Update transaction
	now := time.Now()
	transaction.Status = "failed"
	// Store admin info and reason in Metadata
	transaction.Metadata = fmt.Sprintf(`{"verified_by":%d,"verified_at":"%s","reason":"%s"}`, adminID, now.Format(time.RFC3339), reason)

	if err := s.transactionRepo.Update(transaction); err != nil {
		return err
	}

	// Update order
	order.Status = "cancelled"
	order.PaymentStatus = "failed"
	if err := s.orderRepo.Update(order); err != nil {
		return err
	}

	return nil
}

func (s *orderService) CancelOrder(userID, orderID uint) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// Verify ownership
	if order.UserID != userID {
		return errors.New("unauthorized")
	}

	// Can only cancel pending orders
	if order.PaymentStatus == "paid" || order.Status == "completed" {
		return errors.New("order cannot be cancelled")
	}

	order.Status = "cancelled"
	order.PaymentStatus = "cancelled"
	return s.orderRepo.Update(order)
}
