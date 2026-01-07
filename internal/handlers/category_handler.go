package handlers

import (
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create category (Admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param category body object true "Category data"
// @Success 201 {object} utils.Response
// @Router /categories [post]
// @Security Bearer
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(req.Name, req.Description)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Tags categories
// @Produce json
// @Success 200 {object} utils.Response
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// GetCategoryBySlug godoc
// @Summary Get category by slug
// @Tags categories
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} utils.Response
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.categoryService.GetCategoryBySlug(slug)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// UpdateCategory godoc
// @Summary Update category (Admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body object true "Category data"
// @Success 200 {object} utils.Response
// @Router /categories/{id} [put]
// @Security Bearer
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), req.Name, req.Description)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", category)
}

// DeleteCategory godoc
// @Summary Delete category (Admin only)
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response
// @Router /categories/{id} [delete]
// @Security Bearer
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
