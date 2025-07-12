package handler

import (
	"encoding/json"
	"hot-coffee/help"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"net/http"
)

type OrderHandler struct {
	OrderService *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: service}
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.OrderService.GetAllOrders()
	if err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.OrderService.CreateOrder(order); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var orderStatusUpdate models.Order
	if err := json.NewDecoder(r.Body).Decode(&orderStatusUpdate); err != nil {
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.OrderService.UpdateOrder(orderStatusUpdate); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to update order")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		help.WriteError(w, http.StatusBadRequest, "Missing order_id parameter")
		return
	}
	if err := h.OrderService.DeleteOrder(orderID); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to delete order")
		return
	}
	w.WriteHeader(http.StatusOK)
}
