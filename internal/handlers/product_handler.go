package handlers

import (
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Create product (Admin only)
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product name"
// @Param description formData string false "Product description"
// @Param price formData number true "Product price"
// @Param stock formData int true "Product stock"
// @Param category_id formData int true "Category ID"
// @Param is_featured formData bool false "Is featured"
// @Param is_active formData bool false "Is active"
// @Param image formData file false "Product image"
// @Success 201 {object} utils.Response
// @Router /products [post]
// @Security Bearer
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req services.CreateProductRequest

	// Parse form data
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get file from form
	file, _ := c.FormFile("image")

	// Get user ID from context
	userID := middleware.GetUserID(c)

	product, err := h.productService.CreateProduct(req, file, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", product)
}

// GetAllProducts godoc
// @Summary Get all products
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category_id query int false "Filter by category"
// @Param search query string false "Search by name or description"
// @Success 200 {object} utils.Response
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	var categoryID *uint
	if catID := c.Query("category_id"); catID != "" {
		id, err := strconv.ParseUint(catID, 10, 32)
		if err == nil {
			val := uint(id)
			categoryID = &val
		}
	}

	products, total, err := h.productService.GetAllProducts(page, limit, categoryID, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

// GetProductBySlug godoc
// @Summary Get product by slug
// @Tags products
// @Produce json
// @Param slug path string true "Product slug"
// @Success 200 {object} utils.Response
// @Router /products/slug/{slug} [get]
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	product, err := h.productService.GetProductBySlug(slug)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

// UpdateProduct godoc
// @Summary Update product (Admin only)
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param name formData string false "Product name"
// @Param description formData string false "Product description"
// @Param price formData number false "Product price"
// @Param stock formData int false "Product stock"
// @Param category_id formData int false "Category ID"
// @Param is_featured formData bool false "Is featured"
// @Param is_active formData bool false "Is active"
// @Param image formData file false "Product image"
// @Success 200 {object} utils.Response
// @Router /products/{id} [put]
// @Security Bearer
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req services.UpdateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	file, _ := c.FormFile("image")

	product, err := h.productService.UpdateProduct(uint(id), req, file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary Delete product (Admin only)
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response
// @Router /products/{id} [delete]
// @Security Bearer
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

// GetFeaturedProducts godoc
// @Summary Get featured products
// @Tags products
// @Produce json
// @Param limit query int false "Number of products" default(10)
// @Success 200 {object} utils.Response
// @Router /products/featured [get]
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.productService.GetFeaturedProducts(limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Featured products retrieved successfully", products)
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Tags products
// @Produce json
// @Param category_id path int true "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response
// @Router /products/category/{category_id} [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("category_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productService.GetProductsByCategory(uint(categoryID), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}
