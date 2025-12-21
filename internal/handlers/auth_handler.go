package handlers

import (
	"gin-quickstart/internal/config"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/services"
	"gin-quickstart/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	service services.AuthService
	config  *config.Config
}

func NewAuthHandler(service services.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		config:  config,
	}
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "User data"
// @Success 201 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	authResp, err := h.service.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Registration successful", authResp)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserLoginRequest true "Login credentials"
// @Success 200 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	authResp, err := h.service.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", authResp)
}

// GoogleLogin godoc
// @Summary Google OAuth login
// @Tags auth
// @Produce json
// @Success 302
// @Router /auth/google [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	oauthConfig := &utils.OAuthConfig{
		GoogleClientID:     h.config.GoogleClientID,
		GoogleClientSecret: h.config.GoogleClientSecret,
		GoogleRedirectURL:  h.config.GoogleRedirectURL,
	}

	config := utils.GetGoogleOAuthConfig(oauthConfig)

	// Build URL manually untuk memastikan semua parameter ada
	url := config.AuthCodeURL("state")

	// Return URL dan juga redirect langsung
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code"
// @Success 200 {object} utils.Response
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Authorization code required")
		return
	}

	authResp, err := h.service.GoogleOAuth(code)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Google login successful", authResp)
}

// GithubLogin godoc
// @Summary GitHub OAuth login
// @Tags auth
// @Produce json
// @Success 302
// @Router /auth/github [get]
func (h *AuthHandler) GithubLogin(c *gin.Context) {
	oauthConfig := &utils.OAuthConfig{
		GithubClientID:     h.config.GithubClientID,
		GithubClientSecret: h.config.GithubClientSecret,
		GithubRedirectURL:  h.config.GithubRedirectURL,
	}

	config := utils.GetGithubOAuthConfig(oauthConfig)
	url := config.AuthCodeURL("state", oauth2.SetAuthURLParam("response_type", "code"))

	c.JSON(http.StatusOK, gin.H{
		"auth_url": url,
	})
}

// GithubCallback godoc
// @Summary GitHub OAuth callback
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code"
// @Success 200 {object} utils.Response
// @Router /auth/github/callback [get]
func (h *AuthHandler) GithubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Authorization code required")
		return
	}

	authResp, err := h.service.GithubOAuth(code)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "GitHub login successful", authResp)
}
