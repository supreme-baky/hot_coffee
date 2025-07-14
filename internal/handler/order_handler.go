package handler

import (
	"encoding/json"
	"hot-coffee/help"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
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
		slog.Error("Failed to fetch orders", "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to fetch errors")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		slog.Warn("Invalid order JSON", "error", err)
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.OrderService.CreateOrder(order); err != nil {
		slog.Error("Failed to create order", "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}
	slog.Info("Order created", "orderID", order.ID)
	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request, id string) {
	order, err := h.OrderService.GetOrderByID(id)
	if err != nil {
		slog.Warn("Order not found", "orderID", id)
		help.WriteError(w, http.StatusNotFound, "Order not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request, id string) {
	var updatePayload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updatePayload); err != nil {
		slog.Warn("Invalid status update JSON", "error", err)
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	order, err := h.OrderService.GetOrderByID(id)
	if err != nil {
		slog.Error("Order not found", "orderID", id, "error", err)
		help.WriteError(w, http.StatusNotFound, "Order not found")
		return
	}

	order.Status = updatePayload.Status
	if err := h.OrderService.UpdateOrder(order); err != nil {
		slog.Error("Failed to update order status", "orderID", id, "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	slog.Info("Order status updated", "orderID", id, "status", updatePayload.Status)
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		slog.Warn("Missing order ID in path")
		help.WriteError(w, http.StatusBadRequest, "Missing order ID")
		return
	}
	if err := h.OrderService.DeleteOrder(id); err != nil {
		slog.Error("Failed to delete order", "orderID", id, "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to delete order")
		return
	}
	slog.Info("Order deleted", "orderID", id)
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.OrderService.CloseOrder(id); err != nil {
		slog.Error("Failed to close order", "orderID", id, "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to close order")
		return
	}
	slog.Info("Order closed", "orderID", id)
	w.WriteHeader(http.StatusOK)
}
