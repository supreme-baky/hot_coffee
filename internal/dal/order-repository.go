package dal

import "hot-coffee/models"

type OrderManager interface {
	CreateOrder(order models.Order) error
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id string) (models.Order, error)
	UpdateOrder(order models.Order) error
	DeleteOrder(id string) error
	CloseOrder(id string) error
}
