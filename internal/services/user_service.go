package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error)
	GetProfile(userID uint) (*models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
	UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
	ChangePassword(userID uint, req *models.ChangePasswordRequest) error
	DeleteUser(id uint) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error) {
	// Check if email already exists
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
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	return userResponses, nil
}

func (s *userService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check if new email already used by another user
	if req.Email != user.Email {
		existingUser, _ := s.repo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already used by another user")
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) DeleteUser(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.repo.Delete(user.ID)
}

func (s *userService) GetProfile(userID uint) (*models.UserResponse, error) {
	return s.GetUserByID(userID)
}

func (s *userService) UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	return s.UpdateUser(userID, req)
}

func (s *userService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if user registered with OAuth
	if user.Provider != "local" {
		return errors.New("cannot change password for OAuth users")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.Update(user)
}
