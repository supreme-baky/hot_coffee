package dal

import (
	"encoding/json"
	"errors"
	"hot-coffee/models"
	"log/slog"
	"os"
	"sync"
)

type JSONOrderManager struct {
	filePath string
	orders   []models.Order
	mu       sync.Mutex
}

func NewJSONOrderManager(filePath string) *JSONOrderManager {
	m := &JSONOrderManager{filePath: filePath}
	m.load()
	return m
}

func (m *JSONOrderManager) load() {
	file, err := os.ReadFile(m.filePath)
	if err == nil {
		_ = json.Unmarshal(file, &m.orders)
	}
}

func (m *JSONOrderManager) save() error {
	data, err := json.MarshalIndent(m.orders, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0o644)
}

func (m *JSONOrderManager) CreateOrder(order models.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orders = append(m.orders, order)
	return m.save()
}

func (m *JSONOrderManager) GetAllOrders() ([]models.Order, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.orders, nil
}

func (m *JSONOrderManager) GetOrderByID(id string) (models.Order, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, order := range m.orders {
		if order.ID == id {
			return order, nil
		}
	}
	return models.Order{}, errors.New("order not found")
}

func (m *JSONOrderManager) UpdateOrder(updated models.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, order := range m.orders {
		if order.ID == updated.ID {
			m.orders[i] = updated
			return m.save()
		}
	}
	return errors.New("order not found")
}

func (m *JSONOrderManager) DeleteOrder(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, order := range m.orders {
		if order.ID == id {
			m.orders = append(m.orders[:i], m.orders[i+1:]...)
			return m.save()
		}
	}
	return errors.New("order not found")
}

func (m *JSONOrderManager) CloseOrder(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	slog.Info("Trying to close order", "givenID", id)
	for i, order := range m.orders {
		slog.Info("Checking", "orderID", order.ID)
		if order.ID == id {
			slog.Info("Found and closing", "orderID", id)
			m.orders[i].Status = "closed"
			return m.save()
		}
	}
	slog.Warn("Order not found to close", "givenID", id)
	return errors.New("order not found")
}
