package handlers

import (
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Create order (checkout)
// @Tags orders
// @Accept json
// @Produce json
// @Param order body object true "Order data"
// @Success 201 {object} utils.Response
// @Router /orders [post]
// @Security Bearer
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	order, err := h.orderService.CreateOrder(userID, req.ProductID, req.Quantity)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Order created successfully. Please upload payment proof.", order)
}

// GetUserOrders godoc
// @Summary Get user's orders
// @Tags orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response
// @Router /orders [get]
// @Security Bearer
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.orderService.GetUserOrders(userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Orders retrieved successfully", gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response
// @Router /orders/{id} [get]
// @Security Bearer
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.orderService.GetOrderByID(userID, uint(orderID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order retrieved successfully", order)
}

// UploadPaymentProof godoc
// @Summary Upload payment proof
// @Tags orders
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Order ID"
// @Param proof formData file true "Payment proof image"
// @Success 200 {object} utils.Response
// @Router /orders/{id}/payment-proof [post]
// @Security Bearer
func (h *OrderHandler) UploadPaymentProof(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	file, err := c.FormFile("proof")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payment proof image is required")
		return
	}

	if err := h.orderService.UploadPaymentProof(userID, uint(orderID), file); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment proof uploaded successfully. Waiting for admin verification.", nil)
}

// CancelOrder godoc
// @Summary Cancel order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response
// @Router /orders/{id}/cancel [post]
// @Security Bearer
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	if err := h.orderService.CancelOrder(userID, uint(orderID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order cancelled successfully", nil)
}

// GetAllOrders godoc
// @Summary Get all orders (Admin only)
// @Tags orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} utils.Response
// @Router /admin/orders [get]
// @Security Bearer
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	orders, total, err := h.orderService.GetAllOrders(page, limit, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Orders retrieved successfully", gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// ApprovePayment godoc
// @Summary Approve payment (Admin only)
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response
// @Router /admin/orders/{id}/approve [post]
// @Security Bearer
func (h *OrderHandler) ApprovePayment(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	if err := h.orderService.ApprovePayment(uint(orderID), adminID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment approved successfully", nil)
}

// RejectPayment godoc
// @Summary Reject payment (Admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param reason body object true "Rejection reason"
// @Success 200 {object} utils.Response
// @Router /admin/orders/{id}/reject [post]
// @Security Bearer
func (h *OrderHandler) RejectPayment(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Reason is required")
		return
	}

	if err := h.orderService.RejectPayment(uint(orderID), adminID, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment rejected", nil)
}
