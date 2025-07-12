package dal

import (
	"encoding/json"
	"errors"
	"hot-coffee/models"
	"os"
	"sync"
)

type JSONInventoryManager struct {
	filePath string
	items    []models.InventoryItem
	mu       sync.Mutex
}

func NewJSONInventoryManager(filePath string) *JSONInventoryManager {
	m := &JSONInventoryManager{filePath: filePath}
	m.load()
	return m
}

func (m *JSONInventoryManager) load() {
	file, err := os.ReadFile(m.filePath)
	if err == nil {
		_ = json.Unmarshal(file, &m.items)
	}
}

func (m *JSONInventoryManager) save() error {
	data, err := json.MarshalIndent(m.items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

func (m *JSONInventoryManager) AddNewInventoryItem(item models.InventoryItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = append(m.items, item)
	return m.save()
}

func (m *JSONInventoryManager) GetAllInventoryItems() ([]models.InventoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.items, nil
}

func (m *JSONInventoryManager) GetInventoryItem(id string) (models.InventoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, item := range m.items {
		if item.IngredientID == id {
			return item, nil
		}
	}
	return models.InventoryItem{}, errors.New("item not found")
}

func (m *JSONInventoryManager) UpdateInventoryItem(updated models.InventoryItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, item := range m.items {
		if item.IngredientID == updated.IngredientID {
			m.items[i] = updated
			return m.save()
		}
	}
	return errors.New("item not found")
}

func (m *JSONInventoryManager) DeleteInventoryItem(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, item := range m.items {
		if item.IngredientID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return m.save()
		}
	}
	return errors.New("item not found")
}
