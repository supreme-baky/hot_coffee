package dal

import "hot-coffee/models"

type InventoryManager interface {
	AddNewInventoryItem(item models.InventoryItem) error
	GetAllInventoryItems() ([]models.InventoryItem, error)
	GetInventoryItem(id string) (models.InventoryItem, error)
	UpdateInventoryItem(item models.InventoryItem) error
	DeleteInventoryItem(id string) error
	CheckSufficientIngredients(required []models.MenuItemIngredient) error
	DeductIngredients(required []models.MenuItemIngredient) error
	RestoreIngredients([]models.MenuItemIngredient) error
}
