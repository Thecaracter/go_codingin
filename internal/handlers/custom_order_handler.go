package handlers

import (
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomOrderHandler struct {
	customOrderService services.CustomOrderService
}

func NewCustomOrderHandler(customOrderService services.CustomOrderService) *CustomOrderHandler {
	return &CustomOrderHandler{customOrderService: customOrderService}
}

type CreateCustomOrderRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Requirements string  `json:"requirements"`
	Budget       float64 `json:"budget"`
}

type ProcessCustomOrderRequest struct {
	Status        string   `json:"status" binding:"required"`
	AdminNotes    string   `json:"admin_notes"`
	QuotedPrice   *float64 `json:"quoted_price"`
	EstimatedDays *int     `json:"estimated_days"`
}

// CreateCustomOrder godoc
// @Summary Create custom order request
// @Description User can request custom development work
// @Tags custom-orders
// @Accept json
// @Produce json
// @Param request body CreateCustomOrderRequest true "Custom order data"
// @Success 201 {object} map[string]interface{}
// @Router /custom-orders [post]
// @Security BearerAuth
func (h *CustomOrderHandler) CreateCustomOrder(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateCustomOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.customOrderService.CreateRequest(
		userID,
		req.Title,
		req.Description,
		req.Requirements,
		req.Budget,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Custom order request berhasil dibuat", nil)
}

// GetMyCustomOrders godoc
// @Summary Get user's custom orders
// @Description Get all custom order requests made by the authenticated user
// @Tags custom-orders
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /custom-orders/me [get]
// @Security BearerAuth
func (h *CustomOrderHandler) GetMyCustomOrders(c *gin.Context) {
	userID := c.GetUint("userID")

	customOrders, err := h.customOrderService.GetUserRequests(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Custom orders berhasil diambil", gin.H{
		"custom_orders": customOrders,
	})
}

// GetCustomOrderByID godoc
// @Summary Get custom order by ID
// @Description Get details of a specific custom order
// @Tags custom-orders
// @Produce json
// @Param id path int true "Custom Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /custom-orders/{id} [get]
// @Security BearerAuth
func (h *CustomOrderHandler) GetCustomOrderByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid custom order ID")
		return
	}

	customOrder, err := h.customOrderService.GetRequestByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Custom order berhasil diambil", gin.H{
		"custom_order": customOrder,
	})
}

// CancelCustomOrder godoc
// @Summary Cancel custom order
// @Description User can cancel their own custom order request
// @Tags custom-orders
// @Produce json
// @Param id path int true "Custom Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /custom-orders/{id}/cancel [put]
// @Security BearerAuth
func (h *CustomOrderHandler) CancelCustomOrder(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid custom order ID")
		return
	}

	err = h.customOrderService.CancelRequest(uint(id), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Custom order berhasil dibatalkan", nil)
}

// AdminGetAllCustomOrders godoc
// @Summary Admin get all custom orders
// @Description Admin can view all custom order requests
// @Tags custom-orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{}
// @Router /admin/custom-orders [get]
// @Security BearerAuth
func (h *CustomOrderHandler) AdminGetAllCustomOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	customOrders, total, err := h.customOrderService.GetAllRequests(page, limit, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Custom orders berhasil diambil", gin.H{
		"custom_orders": customOrders,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

// AdminProcessCustomOrder godoc
// @Summary Admin process custom order
// @Description Admin can update status and provide quotation for custom orders
// @Tags custom-orders
// @Accept json
// @Produce json
// @Param id path int true "Custom Order ID"
// @Param request body ProcessCustomOrderRequest true "Process data"
// @Success 200 {object} map[string]interface{}
// @Router /admin/custom-orders/{id}/process [put]
// @Security BearerAuth
func (h *CustomOrderHandler) AdminProcessCustomOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid custom order ID")
		return
	}

	var req ProcessCustomOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.customOrderService.ProcessRequest(
		uint(id),
		req.Status,
		req.AdminNotes,
		req.QuotedPrice,
		req.EstimatedDays,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Custom order berhasil diproses", nil)
}
