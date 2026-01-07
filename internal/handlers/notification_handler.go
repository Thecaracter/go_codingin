package handlers

import (
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// GetMyNotifications godoc
// @Summary Get user's notifications
// @Description Get all notifications for the authenticated user with pagination
// @Tags notifications
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /notifications [get]
// @Security BearerAuth
func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
	userID := c.GetUint("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	notifications, total, err := h.notificationService.GetUserNotifications(userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notifications berhasil diambil", gin.H{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

// GetUnreadNotifications godoc
// @Summary Get unread notifications
// @Description Get all unread notifications for the authenticated user
// @Tags notifications
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /notifications/unread [get]
// @Security BearerAuth
func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.GetUint("userID")

	notifications, err := h.notificationService.GetUnreadNotifications(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Unread notifications berhasil diambil", gin.H{
		"notifications": notifications,
		"count":         len(notifications),
	})
}

// MarkNotificationAsRead godoc
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Router /notifications/{id}/read [put]
// @Security BearerAuth
func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	userID := c.GetUint("userID")
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.notificationService.MarkAsRead(uint(notificationID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification berhasil ditandai sebagai dibaca", nil)
}

// MarkAllNotificationsAsRead godoc
// @Summary Mark all notifications as read
// @Description Mark all user's notifications as read
// @Tags notifications
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /notifications/read-all [put]
// @Security BearerAuth
func (h *NotificationHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	userID := c.GetUint("userID")

	err := h.notificationService.MarkAllAsRead(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Semua notifications berhasil ditandai sebagai dibaca", nil)
}

// DeleteNotification godoc
// @Summary Delete notification
// @Description Delete a specific notification
// @Tags notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Router /notifications/{id} [delete]
// @Security BearerAuth
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID := c.GetUint("userID")
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.notificationService.DeleteNotification(uint(notificationID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification berhasil dihapus", nil)
}
