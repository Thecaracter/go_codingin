package handlers

import (
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	wishlistService services.WishlistService
}

func NewWishlistHandler(wishlistService services.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// AddToWishlist godoc
// @Summary Add item to wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param item body object true "Wishlist item"
// @Success 201 {object} utils.Response
// @Router /wishlist [post]
// @Security Bearer
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	wishlist, err := h.wishlistService.AddToWishlist(userID, req.ProductID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Item added to wishlist", wishlist)
}

// GetUserWishlist godoc
// @Summary Get user's wishlist
// @Tags wishlist
// @Produce json
// @Success 200 {object} utils.Response
// @Router /wishlist [get]
// @Security Bearer
func (h *WishlistHandler) GetUserWishlist(c *gin.Context) {
	userID := middleware.GetUserID(c)

	wishlists, err := h.wishlistService.GetUserWishlist(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Wishlist retrieved successfully", wishlists)
}

// RemoveFromWishlist godoc
// @Summary Remove item from wishlist
// @Tags wishlist
// @Produce json
// @Param id path int true "Wishlist item ID"
// @Success 200 {object} utils.Response
// @Router /wishlist/{id} [delete]
// @Security Bearer
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	wishlistID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid wishlist item ID")
		return
	}

	userID := middleware.GetUserID(c)

	if err := h.wishlistService.RemoveFromWishlist(userID, uint(wishlistID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item removed from wishlist", nil)
}

// ClearWishlist godoc
// @Summary Clear user's wishlist
// @Tags wishlist
// @Produce json
// @Success 200 {object} utils.Response
// @Router /wishlist/clear [delete]
// @Security Bearer
func (h *WishlistHandler) ClearWishlist(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.wishlistService.ClearWishlist(userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Wishlist cleared successfully", nil)
}
