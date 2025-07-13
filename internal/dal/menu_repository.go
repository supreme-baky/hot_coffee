package dal

import "hot-coffee/models"

type MenuManager interface {
	AddNewMenuItem(item models.MenuItem) error
	GetAllMenuItems() ([]models.MenuItem, error)
	GetMenuItem(id string) (models.MenuItem, error)
	UpdateMenuItem(item models.MenuItem) error
	DeleteMenuItem(id string) error
}
func (m *JSONMenuManager) LoadMenuItems() ([]models.MenuItem, error) {
	return m.GetAllMenuItems() // или LoadAll()
}