package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"hot-coffee/models"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"
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
	if err != nil {
		slog.Error("Failed to read inventory file", "path", m.filePath, "error", err)
		return
	}

	if err := json.Unmarshal(file, &m.orders); err != nil {
		slog.Error("Invalid JSON format in inventory file", "path", m.filePath, "error", err)
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

	rand.Seed(time.Now().UnixNano())

	if order.ID == "" {
		for {
			randomID := fmt.Sprintf("%03d", rand.Intn(1000))
			if !m.idExists(randomID) {
				order.ID = randomID
				break
			}
		}
	} else if m.idExists(order.ID) {
		return errors.New("order ID already exists")
	}

	m.orders = append(m.orders, order)
	return m.save()
}

func (m *JSONOrderManager) idExists(id string) bool {
	for _, o := range m.orders {
		if o.ID == id {
			return true
		}
	}
	return false
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

	for i, existing := range m.orders {
		if existing.ID == updated.ID {
			if existing.Status == "closed" {
				return fmt.Errorf("cannot update a closed order (ID: %s)", existing.ID)
			}

			existing.CustomerName = updated.CustomerName
			existing.Items = updated.Items

			m.orders[i] = existing
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
