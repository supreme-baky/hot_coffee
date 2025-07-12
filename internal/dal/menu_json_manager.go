package dal

import (
	"encoding/json"
	"errors"
	"hot-coffee/models"
	"os"
	"sync"
)

type JSONMenuManager struct {
	filePath string
	items    []models.MenuItem
	mu       sync.Mutex
}

func NewJSONMenuManager(filePath string) *JSONMenuManager {
	m := &JSONMenuManager{filePath: filePath}
	m.load()
	return m
}

func (m *JSONMenuManager) load() {
	file, err := os.ReadFile(m.filePath)
	if err == nil {
		_ = json.Unmarshal(file, &m.items)
	}
}

func (m *JSONMenuManager) save() error {
	data, err := json.MarshalIndent(m.items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

func (m *JSONMenuManager) AddNewMenuItem(item models.MenuItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = append(m.items, item)
	return m.save()
}

func (m *JSONMenuManager) GetAllMenuItems() ([]models.MenuItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.items, nil
}

func (m *JSONMenuManager) GetMenuItem(id string) (models.MenuItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return models.MenuItem{}, errors.New("menu item not found")
}

func (m *JSONMenuManager) UpdateMenuItem(updated models.MenuItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, item := range m.items {
		if item.ID == updated.ID {
			m.items[i] = updated
			return m.save()
		}
	}
	return errors.New("menu item not found")
}

func (m *JSONMenuManager) DeleteMenuItem(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, item := range m.items {
		if item.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return m.save()
		}
	}
	return errors.New("menu item not found")
}
