package handlers

import (
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// AddToCart godoc
// @Summary Add item to cart
// @Tags cart
// @Accept json
// @Produce json
// @Param item body object true "Cart item"
// @Success 201 {object} utils.Response
// @Router /cart [post]
// @Security Bearer
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	cart, err := h.cartService.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Item added to cart", cart)
}

// GetUserCart godoc
// @Summary Get user's cart
// @Tags cart
// @Produce json
// @Success 200 {object} utils.Response
// @Router /cart [get]
// @Security Bearer
func (h *CartHandler) GetUserCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	carts, total, err := h.cartService.GetUserCart(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart retrieved successfully", gin.H{
		"items": carts,
		"total": total,
	})
}

// UpdateCartItem godoc
// @Summary Update cart item quantity
// @Tags cart
// @Accept json
// @Produce json
// @Param id path int true "Cart item ID"
// @Param item body object true "Update data"
// @Success 200 {object} utils.Response
// @Router /cart/{id} [put]
// @Security Bearer
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	cartID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	cart, err := h.cartService.UpdateCartItem(userID, uint(cartID), req.Quantity)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart item updated", cart)
}

// RemoveFromCart godoc
// @Summary Remove item from cart
// @Tags cart
// @Produce json
// @Param id path int true "Cart item ID"
// @Success 200 {object} utils.Response
// @Router /cart/{id} [delete]
// @Security Bearer
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	cartID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	userID := middleware.GetUserID(c)

	if err := h.cartService.RemoveFromCart(userID, uint(cartID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item removed from cart", nil)
}

// ClearCart godoc
// @Summary Clear user's cart
// @Tags cart
// @Produce json
// @Success 200 {object} utils.Response
// @Router /cart/clear [delete]
// @Security Bearer
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.cartService.ClearCart(userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart cleared successfully", nil)
}
