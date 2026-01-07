package handlers

import (
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService services.ReviewService
}

func NewReviewHandler(reviewService services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

type CreateReviewRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
}

// CreateReview godoc
// @Summary Create a new review
// @Description User can create a review for a product they purchased
// @Tags reviews
// @Accept json
// @Produce json
// @Param request body CreateReviewRequest true "Review data"
// @Success 201 {object} map[string]interface{}
// @Router /reviews [post]
// @Security BearerAuth
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.reviewService.CreateReview(userID, req.ProductID, req.Rating, req.Comment)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Review berhasil dibuat", nil)
}

// GetProductReviews godoc
// @Summary Get reviews for a product
// @Description Get all reviews for a specific product with pagination
// @Tags reviews
// @Produce json
// @Param product_id path int true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /reviews/product/{product_id} [get]
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.reviewService.GetProductReviews(uint(productID), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	avgRating, _ := h.reviewService.GetAverageRating(uint(productID))

	utils.SuccessResponse(c, http.StatusOK, "Reviews berhasil diambil", gin.H{
		"reviews":        reviews,
		"total":          total,
		"page":           page,
		"limit":          limit,
		"average_rating": avgRating,
	})
}

// GetMyReviews godoc
// @Summary Get user's reviews
// @Description Get all reviews made by the authenticated user
// @Tags reviews
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /reviews/me [get]
// @Security BearerAuth
func (h *ReviewHandler) GetMyReviews(c *gin.Context) {
	userID := c.GetUint("userID")

	reviews, err := h.reviewService.GetUserReviews(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reviews berhasil diambil", gin.H{
		"reviews": reviews,
	})
}

// UpdateReview godoc
// @Summary Update a review
// @Description Update user's own review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param request body UpdateReviewRequest true "Updated review data"
// @Success 200 {object} map[string]interface{}
// @Router /reviews/{id} [put]
// @Security BearerAuth
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := c.GetUint("userID")
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
		return
	}

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.reviewService.UpdateReview(uint(reviewID), userID, req.Rating, req.Comment)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Review berhasil diupdate", nil)
}

// DeleteReview godoc
// @Summary Delete a review
// @Description Delete user's own review
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} map[string]interface{}
// @Router /reviews/{id} [delete]
// @Security BearerAuth
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID := c.GetUint("userID")
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
		return
	}

	err = h.reviewService.DeleteReview(uint(reviewID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Review berhasil dihapus", nil)
}

// AdminDeleteReview godoc
// @Summary Admin delete a review
// @Description Admin can delete any review
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/reviews/{id} [delete]
// @Security BearerAuth
func (h *ReviewHandler) AdminDeleteReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
		return
	}

	err = h.reviewService.AdminDeleteReview(uint(reviewID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Review berhasil dihapus", nil)
}
