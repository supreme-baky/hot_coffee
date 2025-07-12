package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type InventoryService struct {
	InventoryRepo dal.InventoryManager
}

func NewInventoryService(repo dal.InventoryManager) *InventoryService {
	return &InventoryService{
		InventoryRepo: repo,
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

func (s *InventoryService) DeleteInventoryItem(ingredientID string) error {
	return s.InventoryRepo.DeleteInventoryItem(ingredientID)
}

func (s *InventoryService) GetInventoryItem(ingredientID string) (models.InventoryItem, error) {
	return s.InventoryRepo.GetInventoryItem(ingredientID)
}
