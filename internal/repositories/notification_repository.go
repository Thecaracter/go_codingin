package repositories

import (
	"gin-quickstart/internal/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	GetByID(id uint) (*models.Notification, error)
	GetByUserID(userID uint, limit, offset int) ([]*models.Notification, int64, error)
	GetUnreadByUserID(userID uint) ([]*models.Notification, error)
	Update(notification *models.Notification) error
	Delete(id uint) error
	MarkAsRead(id uint) error
	MarkAllAsRead(userID uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) GetByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) GetByUserID(userID uint, limit, offset int) ([]*models.Notification, int64, error) {
	var notifications []*models.Notification
	var total int64

	query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *notificationRepository) GetUnreadByUserID(userID uint) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) Update(notification *models.Notification) error {
	return r.db.Save(notification).Error
}

func (r *notificationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Notification{}, id).Error
}

func (r *notificationRepository) MarkAsRead(id uint) error {
	return r.db.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}
