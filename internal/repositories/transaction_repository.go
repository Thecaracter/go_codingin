package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id uint) (*models.Transaction, error)
	GetByOrderID(orderID uint) (*models.Transaction, error)
	Update(transaction *models.Transaction) error
	GetPendingTransactions(page, limit int) ([]models.Transaction, int64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Order").Preload("Order.User").Preload("Order.Product").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByOrderID(orderID uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Where("order_id = ?", orderID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) GetPendingTransactions(page, limit int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Where("status = ?", "pending")
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Order").Preload("Order.User").Preload("Order.Product").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&transactions).Error

	return transactions, total, err
}
