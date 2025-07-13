package service

import (
	"hot-coffee/internal/dal" // Импортируем dal
	"hot-coffee/models"
)

type OrderService struct {
	OrderRepo dal.OrderManager // Используем dal.OrderManager
}

func NewOrderService(repo dal.OrderManager) *OrderService {
	return &OrderService{
		OrderRepo: repo,
	}
}

func (s *OrderService) CreateOrder(order models.Order) error {
	return s.OrderRepo.CreateOrder(order)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.OrderRepo.GetAllOrders()
}

func (s *OrderService) UpdateOrder(order models.Order) error {
	return s.OrderRepo.UpdateOrder(order)
}

func (s *OrderService) DeleteOrder(orderID string) error {
	return s.OrderRepo.DeleteOrder(orderID)
}

func (s *OrderService) CloseOrder(orderID string) error {
	return s.OrderRepo.CloseOrder(orderID)
}

func (s *OrderService) GetOrderByID(orderID string) (models.Order, error) {
	return s.OrderRepo.GetOrderByID(orderID)
}

