package handler

import (
	"encoding/json"
	"hot-coffee/help"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
)

type MenuHandler struct {
	MenuService *service.MenuService
}

func NewMenuHandler(service *service.MenuService) *MenuHandler {
	return &MenuHandler{MenuService: service}
}

func (h *MenuHandler) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.MenuService.GetAllMenuItems()
	if err != nil {
		slog.Error("Failed to fetch menu items", "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to fetch menu items")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *MenuHandler) AddNewMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Warn("Invalid menu item JSON", "error", err)
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.MenuService.AddNewMenuItem(item); err != nil {
		slog.Error("Failed to add menu item", "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to add menu item")
		return
	}
	slog.Info("Menu item added", "productID", item.ID)
	w.WriteHeader(http.StatusCreated)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.MenuService.GetMenuItem(id)
	if err != nil {
		slog.Warn("Menu item not found", "productID", id)
		help.WriteError(w, http.StatusNotFound, "Menu item not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	var updatedItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		slog.Warn("Invalid JSON for menu item update", "error", err)
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedItem.ID = id
	if err := h.MenuService.UpdateMenuItem(updatedItem); err != nil {
		slog.Error("Failed to update menu item", "productID", id, "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to update menu item")
		return
	}

	slog.Info("Menu item updated", "productID", id)
	w.WriteHeader(http.StatusOK)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		slog.Warn("Missing menu item ID in path")
		help.WriteError(w, http.StatusBadRequest, "Missing menu item ID")
		return
	}

	if err := h.MenuService.DeleteMenuItem(id); err != nil {
		slog.Error("Failed to delete menu item", "productID", id, "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to delete menu item")
		return
	}

	slog.Info("Menu item deleted", "productID", id)
	w.WriteHeader(http.StatusOK)
}
