package service

import (
	"errors"
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
	menuItems, err := s.MenuRepo.GetAllMenuItems()
	if err != nil {
		return err
	}

	menuMap := make(map[string]models.MenuItem)
	for _, item := range menuItems {
		menuMap[item.ID] = item
	}

	var ingredientsList []models.MenuItemIngredient
	for _, orderItem := range order.Items {
		menuItem, ok := menuMap[orderItem.ProductID]
		if !ok {
			return errors.New("invalid product ID: " + orderItem.ProductID)
		}
		for _, ing := range menuItem.Ingredients {
			ingredientsList = append(ingredientsList, models.MenuItemIngredient{
				IngredientID: ing.IngredientID,
				Quantity:     ing.Quantity * float64(orderItem.Quantity),
			})
		}
	}

	if err := s.InventoryRepo.CheckSufficientIngredients(ingredientsList); err != nil {
		return err
	}
	if err := s.InventoryRepo.DeductIngredients(ingredientsList); err != nil {
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
