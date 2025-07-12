package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type MenuService struct {
	MenuRepo dal.MenuManager
}

func NewMenuService(repo dal.MenuManager) *MenuService {
	return &MenuService{
		MenuRepo: repo,
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

func (s *MenuService) DeleteMenuItem(menuItemID string) error {
	return s.MenuRepo.DeleteMenuItem(menuItemID)
}
