package handlers

import (
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	downloadService services.DownloadService
}

func NewDownloadHandler(downloadService services.DownloadService) *DownloadHandler {
	return &DownloadHandler{
		downloadService: downloadService,
	}
}

// DownloadProduct godoc
// @Summary Download product file
// @Tags downloads
// @Produce json
// @Param product_id query int true "Product ID"
// @Param order_id query int true "Order ID"
// @Success 200 {object} utils.Response
// @Router /downloads [post]
// @Security Bearer
func (h *DownloadHandler) DownloadProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)

	productID, err := strconv.ParseUint(c.Query("product_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	orderID, err := strconv.ParseUint(c.Query("order_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// Check if user can download
	canDownload, err := h.downloadService.CanDownload(userID, uint(productID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !canDownload {
		utils.ErrorResponse(c, http.StatusForbidden, "You don't have access to download this product")
		return
	}

	// Create download record
	if err := h.downloadService.CreateDownload(userID, uint(orderID), uint(productID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Download recorded successfully", gin.H{
		"message": "Product download will start shortly",
		"note":    "In production, this would return a secure download URL",
	})
}

// GetUserDownloads godoc
// @Summary Get user's download history
// @Tags downloads
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response
// @Router /downloads [get]
// @Security Bearer
func (h *DownloadHandler) GetUserDownloads(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	downloads, total, err := h.downloadService.GetUserDownloads(userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Downloads retrieved successfully", gin.H{
		"downloads": downloads,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// GetDownloadHistory godoc
// @Summary Get download history for specific product
// @Tags downloads
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} utils.Response
// @Router /downloads/history/{product_id} [get]
// @Security Bearer
func (h *DownloadHandler) GetDownloadHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	downloads, err := h.downloadService.GetDownloadHistory(userID, uint(productID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Download history retrieved successfully", downloads)
}
