package service

import (
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type InventoryService struct {
	InventoryRepo dal.InventoryManager
	MenuRepo      dal.MenuManager
}

func NewInventoryService(invRepo dal.InventoryManager, menuRepo dal.MenuManager) *InventoryService {
	return &InventoryService{
		InventoryRepo: invRepo,
		MenuRepo:      menuRepo,
	}
}

func (s *InventoryService) GetAllInventoryItems() ([]models.InventoryItem, error) {
	return s.InventoryRepo.GetAllInventoryItems()
}

func (s *InventoryService) AddNewInventoryItem(item models.InventoryItem) error {
	return s.InventoryRepo.AddNewInventoryItem(item)
}

func (s *InventoryService) UpdateInventoryItem(item models.InventoryItem) error {
	return s.InventoryRepo.UpdateInventoryItem(item)
}

func (s *InventoryService) DeleteInventoryItem(id string) error {
	menuItems, err := s.MenuRepo.GetAllMenuItems()
	if err != nil {
		return fmt.Errorf("failed to load menu items: %w", err)
	}

	for _, menuItem := range menuItems {
		for _, ingredient := range menuItem.Ingredients {
			if ingredient.IngredientID == id {
				return fmt.Errorf("cannot delete inventory item '%s': used in menu item '%s'", id, menuItem.Name)
			}
		}
	}

	return s.InventoryRepo.DeleteInventoryItem(id)
}

func (s *InventoryService) GetInventoryItem(ingredientID string) (models.InventoryItem, error) {
	return s.InventoryRepo.GetInventoryItem(ingredientID)
}
