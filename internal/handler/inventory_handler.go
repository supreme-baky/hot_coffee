package handler

import (
	"encoding/json"
	"hot-coffee/help"
	"hot-coffee/internal/service"
	"hot-coffee/models"
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
		help.WriteError(w, http.StatusInternalServerError, "Failed to fetch inventory items")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *InventoryHandler) AddNewInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.InventoryService.AddNewInventoryItem(item); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to add inventory item")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.InventoryService.UpdateInventoryItem(item); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to update inventory item")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	ingredientID := r.URL.Query().Get("ingredient_id")
	if ingredientID == "" {
		help.WriteError(w, http.StatusBadRequest, "Missing ingredient_id parameter")
		return
	}
	if err := h.InventoryService.DeleteInventoryItem(ingredientID); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to delete inventory item")
		return
	}
	w.WriteHeader(http.StatusOK)
}
