package handler

import (
	"encoding/json"
	"hot-coffee/help"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
)

type InventoryHandler struct {
	InventoryService *service.InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{InventoryService: service}
}

func (h *InventoryHandler) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.InventoryService.GetAllInventoryItems()
	if err != nil {
		slog.Error("Failed to fetch inventory items", "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to fetch inventory items")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *InventoryHandler) AddNewInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Warn("Invalid inventory JSON", "error", err)
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.InventoryService.AddNewInventoryItem(item); err != nil {
		slog.Error("Failed to add inventory item", "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to add inventory item")
		return
	}
	slog.Info("Inventory item added", "ingredientID", item.IngredientID)
	w.WriteHeader(http.StatusCreated)
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.InventoryService.GetInventoryItem(id)
	if err != nil {
		slog.Warn("Inventory item not found", "ingredientID", id)
		help.WriteError(w, http.StatusNotFound, "Inventory item not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	var updatedItem models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		slog.Warn("Invalid JSON for inventory update", "error", err)
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedItem.IngredientID = id
	if err := h.InventoryService.UpdateInventoryItem(updatedItem); err != nil {
		slog.Error("Failed to update inventory item", "ingredientID", id, "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to update inventory item")
		return
	}

	slog.Info("Inventory item updated", "ingredientID", id)
	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		slog.Warn("Missing inventory item ID in path")
		help.WriteError(w, http.StatusBadRequest, "Missing inventory item ID")
		return
	}

	if err := h.InventoryService.DeleteInventoryItem(id); err != nil {
		slog.Error("Failed to delete inventory item", "ingredientID", id, "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to delete inventory item")
		return
	}

	slog.Info("Inventory item deleted", "ingredientID", id)
	w.WriteHeader(http.StatusOK)
}
