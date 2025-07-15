package service

import (
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type MenuService struct {
	MenuRepo  dal.MenuManager
	OrderRepo dal.OrderManager
}

func NewMenuService(menuRepo dal.MenuManager, orderRepo dal.OrderManager) *MenuService {
	return &MenuService{
		MenuRepo:  menuRepo,
		OrderRepo: orderRepo,
	}
}

func (s *MenuService) AddNewMenuItem(item models.MenuItem) error {
	return s.MenuRepo.AddNewMenuItem(item)
}

func (s *MenuService) GetAllMenuItems() ([]models.MenuItem, error) {
	return s.MenuRepo.GetAllMenuItems()
}

func (s *MenuService) GetMenuItem(menuItemID string) (models.MenuItem, error) {
	return s.MenuRepo.GetMenuItem(menuItemID)
}

func (s *MenuService) UpdateMenuItem(item models.MenuItem) error {
	return s.MenuRepo.UpdateMenuItem(item)
}

func (s *MenuService) DeleteMenuItem(id string) error {
	orders, err := s.OrderRepo.GetAllOrders()
	if err != nil {
		return fmt.Errorf("failed to check orders: %w", err)
	}

	for _, order := range orders {
		if order.Status == "open" {
			for _, item := range order.Items {
				if item.ProductID == id {
					return fmt.Errorf("cannot delete menu item '%s': it is used in open order '%s'", id, order.ID)
				}
			}
		}
	}

	return s.MenuRepo.DeleteMenuItem(id)
}
