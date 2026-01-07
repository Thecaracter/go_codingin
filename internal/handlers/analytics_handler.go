package handlers

import (
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Description Get overall dashboard statistics for admin
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/analytics/dashboard [get]
// @Security BearerAuth
func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.analyticsService.GetDashboardStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard stats berhasil diambil", gin.H{
		"stats": stats,
	})
}

// GetRevenueStats godoc
// @Summary Get revenue statistics
// @Description Get revenue statistics for a date range
// @Tags analytics
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)" default(30 days ago)
// @Param end_date query string false "End date (YYYY-MM-DD)" default(today)
// @Success 200 {object} map[string]interface{}
// @Router /admin/analytics/revenue [get]
// @Security BearerAuth
func (h *AnalyticsHandler) GetRevenueStats(c *gin.Context) {
	// Default to last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// Parse dates if provided
	if startStr := c.Query("start_date"); startStr != "" {
		if parsed, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = parsed
		}
	}

	if endStr := c.Query("end_date"); endStr != "" {
		if parsed, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = parsed
		}
	}

	stats, err := h.analyticsService.GetRevenueStats(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Revenue stats berhasil diambil", gin.H{
		"stats": stats,
	})
}

// GetTopProducts godoc
// @Summary Get top selling products
// @Description Get list of top selling products
// @Tags analytics
// @Produce json
// @Param limit query int false "Number of products to return" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /admin/analytics/top-products [get]
// @Security BearerAuth
func (h *AnalyticsHandler) GetTopProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.analyticsService.GetTopProducts(limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Top products berhasil diambil", gin.H{
		"products": products,
	})
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get user statistics including registrations and roles
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/analytics/users [get]
// @Security BearerAuth
func (h *AnalyticsHandler) GetUserStats(c *gin.Context) {
	stats, err := h.analyticsService.GetUserStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User stats berhasil diambil", gin.H{
		"stats": stats,
	})
}

// GetOrderStats godoc
// @Summary Get order statistics
// @Description Get order statistics by status and payment status
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/analytics/orders [get]
// @Security BearerAuth
func (h *AnalyticsHandler) GetOrderStats(c *gin.Context) {
	stats, err := h.analyticsService.GetOrderStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order stats berhasil diambil", gin.H{
		"stats": stats,
	})
}
