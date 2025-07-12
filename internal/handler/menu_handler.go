package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"net/http"
)

type MenuHandler struct {
	MenuService service.MenuService
}

func NewMenuHandler(service service.MenuService) *MenuHandler {
	return &MenuHandler{MenuService: service}
}

func (h *MenuHandler) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.MenuService.GetAllMenuItems()
	if err != nil {
		http.Error(w, "Failed to fetch menu items", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *MenuHandler) AddNewMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.MenuService.AddNewMenuItem(item); err != nil {
		http.Error(w, "Failed to add menu item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.MenuService.UpdateMenuItem(item); err != nil {
		http.Error(w, "Failed to update menu item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	menuItemID := r.URL.Query().Get("menu_item_id")
	if menuItemID == "" {
		http.Error(w, "Missing menu_item_id parameter", http.StatusBadRequest)
		return
	}
	if err := h.MenuService.DeleteMenuItem(menuItemID); err != nil {
		http.Error(w, "Failed to delete menu item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
