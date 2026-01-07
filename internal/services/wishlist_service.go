package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
)

type WishlistService interface {
	AddToWishlist(userID, productID uint) (*models.Wishlist, error)
	GetUserWishlist(userID uint) ([]models.Wishlist, error)
	RemoveFromWishlist(userID, wishlistID uint) error
	ClearWishlist(userID uint) error
	IsInWishlist(userID, productID uint) (bool, error)
}

type wishlistService struct {
	wishlistRepo repositories.WishlistRepository
	productRepo  repositories.ProductRepository
}

func NewWishlistService(wishlistRepo repositories.WishlistRepository, productRepo repositories.ProductRepository) WishlistService {
	return &wishlistService{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
	}
}

func (s *wishlistService) AddToWishlist(userID, productID uint) (*models.Wishlist, error) {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if !product.IsActive {
		return nil, errors.New("product is not available")
	}

	existing, err := s.wishlistRepo.GetByUserAndProduct(userID, productID)
	if err == nil && existing != nil {
		return nil, errors.New("product already in wishlist")
	}

	wishlist := &models.Wishlist{
		UserID:    userID,
		ProductID: productID,
	}

	if err := s.wishlistRepo.Create(wishlist); err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (s *wishlistService) GetUserWishlist(userID uint) ([]models.Wishlist, error) {
	return s.wishlistRepo.GetByUserID(userID)
}

func (s *wishlistService) RemoveFromWishlist(userID, wishlistID uint) error {
	wishlist, err := s.wishlistRepo.GetByID(wishlistID)
	if err != nil {
		return errors.New("wishlist item not found")
	}

	if wishlist.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.wishlistRepo.Delete(wishlistID)
}

func (s *wishlistService) ClearWishlist(userID uint) error {
	return s.wishlistRepo.DeleteByUserID(userID)
}

func (s *wishlistService) IsInWishlist(userID, productID uint) (bool, error) {
	wishlist, err := s.wishlistRepo.GetByUserAndProduct(userID, productID)
	if err != nil {
		return false, nil
	}
	return wishlist != nil, nil
}
