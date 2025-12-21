package services

import (
	"context"
	"errors"
	"gin-quickstart/internal/config"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"gin-quickstart/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *models.UserCreateRequest) (*models.AuthResponse, error)
	Login(req *models.UserLoginRequest) (*models.AuthResponse, error)
	GoogleOAuth(code string) (*models.AuthResponse, error)
	GithubOAuth(code string) (*models.AuthResponse, error)
}

type authService struct {
	repo   repositories.UserRepository
	config *config.Config
}

func NewAuthService(repo repositories.UserRepository, config *config.Config) AuthService {
	return &authService{
		repo:   repo,
		config: config,
	}
}

func (s *authService) Register(req *models.UserCreateRequest) (*models.AuthResponse, error) {
	// Check if email exists
	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Role:       "user",
		Provider:   "local",
		IsVerified: false,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User: models.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
			Provider:   user.Provider,
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt,
		},
		Token: token,
	}, nil
}

func (s *authService) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if user registered with OAuth
	if user.Provider != "local" {
		return nil, errors.New("please login with " + user.Provider)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User: models.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
			Provider:   user.Provider,
			AvatarURL:  user.AvatarURL,
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt,
		},
		Token: token,
	}, nil
}

func (s *authService) GoogleOAuth(code string) (*models.AuthResponse, error) {
	oauthConfig := &utils.OAuthConfig{
		GoogleClientID:     s.config.GoogleClientID,
		GoogleClientSecret: s.config.GoogleClientSecret,
		GoogleRedirectURL:  s.config.GoogleRedirectURL,
	}

	googleConfig := utils.GetGoogleOAuthConfig(oauthConfig)
	userInfo, err := utils.GetGoogleUserInfo(context.Background(), code, googleConfig)
	if err != nil {
		return nil, err
	}

	return s.handleOAuthLogin(userInfo)
}

func (s *authService) GithubOAuth(code string) (*models.AuthResponse, error) {
	oauthConfig := &utils.OAuthConfig{
		GithubClientID:     s.config.GithubClientID,
		GithubClientSecret: s.config.GithubClientSecret,
		GithubRedirectURL:  s.config.GithubRedirectURL,
	}

	githubConfig := utils.GetGithubOAuthConfig(oauthConfig)
	userInfo, err := utils.GetGithubUserInfo(context.Background(), code, githubConfig)
	if err != nil {
		return nil, err
	}

	return s.handleOAuthLogin(userInfo)
}

func (s *authService) handleOAuthLogin(userInfo *utils.OAuthUserInfo) (*models.AuthResponse, error) {
	// Check if user exists
	user, err := s.repo.FindByEmail(userInfo.Email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new user if not exists
	if user == nil {
		user = &models.User{
			Email:      userInfo.Email,
			Name:       userInfo.Name,
			Provider:   userInfo.Provider,
			ProviderID: userInfo.ProviderID,
			AvatarURL:  userInfo.AvatarURL,
			Role:       "user",
			IsVerified: true, // OAuth users are auto-verified
		}

		if err := s.repo.Create(user); err != nil {
			return nil, err
		}
	} else {
		// Update user info if exists
		user.Name = userInfo.Name
		user.AvatarURL = userInfo.AvatarURL
		user.IsVerified = true
		if err := s.repo.Update(user); err != nil {
			return nil, err
		}
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User: models.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
			Provider:   user.Provider,
			AvatarURL:  user.AvatarURL,
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt,
		},
		Token: token,
	}, nil
}
