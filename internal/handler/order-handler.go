package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"net/http"
)

type OrderHandler struct {
	OrderService service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: service}
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.OrderService.GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.OrderService.CreateOrder(order); err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var orderStatusUpdate models.Order
	if err := json.NewDecoder(r.Body).Decode(&orderStatusUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.OrderService.UpdateOrder(orderStatusUpdate); err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(w, "Missing order_id parameter", http.StatusBadRequest)
		return
	}
	if err := h.OrderService.DeleteOrder(orderID); err != nil {
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
