package handlers

import (
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetProfile godoc
// @Summary Get user profile
// @Security Bearer
// @Tags user
// @Produce json
// @Success 200 {object} utils.Response
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.service.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Security Bearer
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UserUpdateRequest true "User data"
// @Success 200 {object} utils.Response
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", user)
}

// ChangePassword godoc
// @Summary Change password
// @Security Bearer
// @Tags user
// @Accept json
// @Produce json
// @Param password body models.ChangePasswordRequest true "Password data"
// @Success 200 {object} utils.Response
// @Router /user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.service.ChangePassword(userID, &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// DeleteAccount godoc
// @Summary Delete user account
// @Security Bearer
// @Tags user
// @Produce json
// @Success 200 {object} utils.Response
// @Router /user/account [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.service.DeleteUser(userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account deleted successfully", nil)
}

// Admin: GetAllUsers godoc
// @Summary Get all users (Admin)
// @Security Bearer
// @Tags admin
// @Produce json
// @Success 200 {object} utils.Response
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", users)
}

// Admin: GetUserByID godoc
// @Summary Get user by ID (Admin)
// @Security Bearer
// @Tags admin
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

// Admin: UpdateUser godoc
// @Summary Update user (Admin)
// @Security Bearer
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UserUpdateRequest true "User data"
// @Success 200 {object} utils.Response
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.service.UpdateUser(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", user)
}

// Admin: DeleteUser godoc
// @Summary Delete user (Admin)
// @Security Bearer
// @Tags admin
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// CreateUser is deprecated - use Register endpoint instead
func (h *UserHandler) CreateUser(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusGone, "This endpoint is deprecated. Use /api/v1/auth/register instead")
}
