package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type OrderService struct {
	OrderRepo dal.OrderManager
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

func (s *OrderService) UpdateOrder(orderID models.Order) error {
	return s.OrderRepo.UpdateOrder(orderID)
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
