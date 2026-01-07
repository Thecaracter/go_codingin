package services

import (
	"errors"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
)

type CartService interface {
	AddToCart(userID, productID uint, quantity int) (*models.Cart, error)
	GetUserCart(userID uint) ([]models.Cart, float64, error)
	UpdateCartItem(userID, cartID uint, quantity int) (*models.Cart, error)
	RemoveFromCart(userID, cartID uint) error
	ClearCart(userID uint) error
}

type cartService struct {
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository
}

func NewCartService(cartRepo repositories.CartRepository, productRepo repositories.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) AddToCart(userID, productID uint, quantity int) (*models.Cart, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if !product.IsActive {
		return nil, errors.New("product is not available")
	}

	existingCart, err := s.cartRepo.GetByUserAndProduct(userID, productID)
	if err == nil {
		newQuantity := existingCart.Quantity + quantity
		existingCart.Quantity = newQuantity
		if err := s.cartRepo.Update(existingCart); err != nil {
			return nil, err
		}
		return existingCart, nil
	}

	cart := &models.Cart{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	if err := s.cartRepo.Create(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) GetUserCart(userID uint) ([]models.Cart, float64, error) {
	carts, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	var total float64
	for _, cart := range carts {
		total += cart.Product.Price * float64(cart.Quantity)
	}

	return carts, total, nil
}

func (s *cartService) UpdateCartItem(userID, cartID uint, quantity int) (*models.Cart, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	cart, err := s.cartRepo.GetByID(cartID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	if cart.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	_, err = s.productRepo.GetByID(cart.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	cart.Quantity = quantity
	if err := s.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) RemoveFromCart(userID, cartID uint) error {
	cart, err := s.cartRepo.GetByID(cartID)
	if err != nil {
		return errors.New("cart item not found")
	}

	if cart.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.cartRepo.Delete(cartID)
}

func (s *cartService) ClearCart(userID uint) error {
	return s.cartRepo.DeleteByUserID(userID)
}
