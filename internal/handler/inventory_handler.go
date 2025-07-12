package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"net/http"
)

type InventoryHandler struct {
	InventoryService service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{InventoryService: service}
}

func (h *InventoryHandler) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.InventoryService.GetAllInventoryItems()
	if err != nil {
		http.Error(w, "Failed to fetch inventory items", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *InventoryHandler) AddNewInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.InventoryService.AddNewInventoryItem(item); err != nil {
		http.Error(w, "Failed to add inventory item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.InventoryService.UpdateInventoryItem(item); err != nil {
		http.Error(w, "Failed to update inventory item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	ingredientID := r.URL.Query().Get("ingredient_id")
	if ingredientID == "" {
		http.Error(w, "Missing ingredient_id parameter", http.StatusBadRequest)
		return
	}
	if err := h.InventoryService.DeleteInventoryItem(ingredientID); err != nil {
		http.Error(w, "Failed to delete inventory item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
