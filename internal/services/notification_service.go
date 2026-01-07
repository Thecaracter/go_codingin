package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
)

type NotificationService interface {
	CreateNotification(userID uint, notifType, title, message string) error
	GetUserNotifications(userID uint, page, limit int) ([]*models.Notification, int64, error)
	GetUnreadNotifications(userID uint) ([]*models.Notification, error)
	MarkAsRead(notificationID, userID uint) error
	MarkAllAsRead(userID uint) error
	DeleteNotification(notificationID, userID uint) error
}

type notificationService struct {
	notificationRepo repositories.NotificationRepository
}

func NewNotificationService(notificationRepo repositories.NotificationRepository) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *notificationService) CreateNotification(userID uint, notifType, title, message string) error {
	if title == "" || message == "" {
		return errors.New("title dan message harus diisi")
	}

	// Valid notification types
	validTypes := map[string]bool{
		"order":        true,
		"payment":      true,
		"download":     true,
		"review":       true,
		"custom_order": true,
		"system":       true,
	}

	if !validTypes[notifType] {
		notifType = "system"
	}

	notification := &models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		IsRead:  false,
	}

	return s.notificationRepo.Create(notification)
}

func (s *notificationService) GetUserNotifications(userID uint, page, limit int) ([]*models.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.notificationRepo.GetByUserID(userID, limit, offset)
}

func (s *notificationService) GetUnreadNotifications(userID uint) ([]*models.Notification, error) {
	return s.notificationRepo.GetUnreadByUserID(userID)
}

func (s *notificationService) MarkAsRead(notificationID, userID uint) error {
	notification, err := s.notificationRepo.GetByID(notificationID)
	if err != nil {
		return errors.New("notification tidak ditemukan")
	}

	if notification.UserID != userID {
		return errors.New("anda tidak memiliki akses untuk notification ini")
	}

	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *notificationService) MarkAllAsRead(userID uint) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *notificationService) DeleteNotification(notificationID, userID uint) error {
	notification, err := s.notificationRepo.GetByID(notificationID)
	if err != nil {
		return errors.New("notification tidak ditemukan")
	}

	if notification.UserID != userID {
		return errors.New("anda tidak memiliki akses untuk menghapus notification ini")
	}

	return s.notificationRepo.Delete(notificationID)
}
