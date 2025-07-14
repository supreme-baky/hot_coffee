package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"hot-coffee/models"
	"log/slog"
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
	if err != nil {
		slog.Error("Failed to read inventory file", "path", m.filePath, "error", err)
		return
	}

	if err := json.Unmarshal(file, &m.items); err != nil {
		slog.Error("Invalid JSON format in inventory file", "path", m.filePath, "error", err)
	}
}

func (m *JSONInventoryManager) save() error {
	data, err := json.MarshalIndent(m.items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0o644)
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

func (m *JSONInventoryManager) CheckSufficientIngredients(required []models.MenuItemIngredient) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, req := range required {
		found := false
		for _, inv := range m.items {
			if inv.IngredientID == req.IngredientID {
				found = true
				if inv.Quantity < req.Quantity {
					return fmt.Errorf(
						"insufficient inventory for ingredient '%s'. Required: %.2f, Available: %.2f",
						inv.Name, req.Quantity, inv.Quantity,
					)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("ingredient '%s' not found in inventory", req.IngredientID)
		}
	}
	return nil
}

func (m *JSONInventoryManager) DeductIngredients(required []models.MenuItemIngredient) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, req := range required {
		for i := range m.items {
			if m.items[i].IngredientID == req.IngredientID {
				m.items[i].Quantity -= req.Quantity
				break
			}
		}
	}
	return m.save()
}
