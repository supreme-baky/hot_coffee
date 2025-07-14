package service

import (
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"time"
)

type OrderService struct {
	OrderRepo     dal.OrderManager
	MenuRepo      dal.MenuManager
	InventoryRepo dal.InventoryManager
}

func NewOrderService(orderRepo dal.OrderManager, menuRepo dal.MenuManager, inventoryRepo dal.InventoryManager) *OrderService {
	return &OrderService{
		OrderRepo:     orderRepo,
		MenuRepo:      menuRepo,
		InventoryRepo: inventoryRepo,
	}
}

func (s *OrderService) CreateOrder(order models.Order) error {
	requiredIngredients := make(map[string]float64)

	for _, item := range order.Items {
		menuItem, err := s.MenuRepo.GetMenuItem(item.ProductID)
		if err != nil {
			return fmt.Errorf("Invalid product ID '%s'", item.ProductID)
		}
		for _, ing := range menuItem.Ingredients {
			requiredIngredients[ing.IngredientID] += ing.Quantity * float64(item.Quantity)
		}
	}

	var requiredList []models.MenuItemIngredient
	for id, qty := range requiredIngredients {
		requiredList = append(requiredList, models.MenuItemIngredient{
			IngredientID: id,
			Quantity:     qty,
		})
	}

	if err := s.InventoryRepo.CheckSufficientIngredients(requiredList); err != nil {
		return err
	}

	if err := s.InventoryRepo.DeductIngredients(requiredList); err != nil {
		return err
	}

	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

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
