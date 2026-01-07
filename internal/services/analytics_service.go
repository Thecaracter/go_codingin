package services

import (
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"time"
)

type AnalyticsService interface {
	GetDashboardStats() (map[string]interface{}, error)
	GetRevenueStats(startDate, endDate time.Time) (map[string]interface{}, error)
	GetTopProducts(limit int) ([]*models.Product, error)
	GetUserStats() (map[string]interface{}, error)
	GetOrderStats() (map[string]interface{}, error)
}

type analyticsService struct {
	userRepo    repositories.UserRepository
	productRepo repositories.ProductRepository
	orderRepo   repositories.OrderRepository
	reviewRepo  repositories.ReviewRepository
}

func NewAnalyticsService(
	userRepo repositories.UserRepository,
	productRepo repositories.ProductRepository,
	orderRepo repositories.OrderRepository,
	reviewRepo repositories.ReviewRepository,
) AnalyticsService {
	return &analyticsService{
		userRepo:    userRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
		reviewRepo:  reviewRepo,
	}
}

func (s *analyticsService) GetDashboardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total users
	users, _, _ := s.userRepo.GetAllUsers(1, 100000)
	stats["total_users"] = len(users)

	// Total products
	products, _, _ := s.productRepo.GetAll(1, 100000, nil, "")
	stats["total_products"] = len(products)

	// Total orders
	orders, totalOrders, _ := s.orderRepo.GetAll(1, 100000, "")
	stats["total_orders"] = totalOrders

	// Calculate total revenue
	var totalRevenue float64
	for _, order := range orders {
		if order.Status == "completed" {
			totalRevenue += order.TotalAmount
		}
	}
	stats["total_revenue"] = totalRevenue

	// Orders by status
	ordersByStatus := make(map[string]int)
	for _, order := range orders {
		ordersByStatus[order.Status]++
	}
	stats["orders_by_status"] = ordersByStatus

	// Recent orders (last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	recentOrders := 0
	for _, order := range orders {
		if order.CreatedAt.After(sevenDaysAgo) {
			recentOrders++
		}
	}
	stats["recent_orders"] = recentOrders

	return stats, nil
}

func (s *analyticsService) GetRevenueStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	orders, _, err := s.orderRepo.GetAll(1, 100000, "")
	if err != nil {
		return nil, err
	}

	var totalRevenue float64
	var completedOrders int
	revenueByDate := make(map[string]float64)

	for _, order := range orders {
		if order.Status == "completed" &&
			order.CreatedAt.After(startDate) &&
			order.CreatedAt.Before(endDate) {

			totalRevenue += order.TotalAmount
			completedOrders++

			// Group by date
			dateKey := order.CreatedAt.Format("2006-01-02")
			revenueByDate[dateKey] += order.TotalAmount
		}
	}

	stats["total_revenue"] = totalRevenue
	stats["completed_orders"] = completedOrders
	stats["revenue_by_date"] = revenueByDate
	stats["start_date"] = startDate.Format("2006-01-02")
	stats["end_date"] = endDate.Format("2006-01-02")

	if completedOrders > 0 {
		stats["average_order_value"] = totalRevenue / float64(completedOrders)
	} else {
		stats["average_order_value"] = 0
	}

	return stats, nil
}

func (s *analyticsService) GetTopProducts(limit int) ([]*models.Product, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get all orders
	orders, _, err := s.orderRepo.GetAll(1, 100000, "completed")
	if err != nil {
		return nil, err
	}

	// Count sales per product
	productSales := make(map[uint]int)
	for _, order := range orders {
		if order.ProductID != nil {
			productSales[*order.ProductID]++
		}
	}

	// Get products and sort by sales
	products, _, err := s.productRepo.GetAll(1, 100000, nil, "")
	if err != nil {
		return nil, err
	}

	// Sort products by sales count
	type productWithSales struct {
		product *models.Product
		sales   int
	}

	var productsWithSales []productWithSales
	for i := range products {
		productsWithSales = append(productsWithSales, productWithSales{
			product: &products[i],
			sales:   productSales[products[i].ID],
		})
	}

	// Simple bubble sort by sales (descending)
	for i := 0; i < len(productsWithSales); i++ {
		for j := i + 1; j < len(productsWithSales); j++ {
			if productsWithSales[j].sales > productsWithSales[i].sales {
				productsWithSales[i], productsWithSales[j] = productsWithSales[j], productsWithSales[i]
			}
		}
	}

	// Get top products
	var topProducts []*models.Product
	for i := 0; i < limit && i < len(productsWithSales); i++ {
		topProducts = append(topProducts, productsWithSales[i].product)
	}

	return topProducts, nil
}

func (s *analyticsService) GetUserStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	users, _, err := s.userRepo.GetAllUsers(1, 100000)
	if err != nil {
		return nil, err
	}

	stats["total_users"] = len(users)

	// Count by role
	usersByRole := make(map[string]int)
	for _, user := range users {
		usersByRole[user.Role]++
	}
	stats["users_by_role"] = usersByRole

	// New users (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	newUsers := 0
	for _, user := range users {
		if user.CreatedAt.After(thirtyDaysAgo) {
			newUsers++
		}
	}
	stats["new_users_30_days"] = newUsers

	return stats, nil
}

func (s *analyticsService) GetOrderStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	orders, totalOrders, err := s.orderRepo.GetAll(1, 100000, "")
	if err != nil {
		return nil, err
	}

	stats["total_orders"] = totalOrders

	// Orders by status
	ordersByStatus := make(map[string]int)
	for _, order := range orders {
		ordersByStatus[order.Status]++
	}
	stats["orders_by_status"] = ordersByStatus

	// Orders by payment status
	ordersByPaymentStatus := make(map[string]int)
	for _, order := range orders {
		ordersByPaymentStatus[order.PaymentStatus]++
	}
	stats["orders_by_payment_status"] = ordersByPaymentStatus

	// Calculate conversion rate (completed / total)
	if totalOrders > 0 {
		completedOrders := ordersByStatus["completed"]
		stats["conversion_rate"] = float64(completedOrders) / float64(totalOrders) * 100
	} else {
		stats["conversion_rate"] = 0
	}

	return stats, nil
}
