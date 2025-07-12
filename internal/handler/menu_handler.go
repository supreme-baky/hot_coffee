package handler

import (
	"encoding/json"
	"hot-coffee/help"
	"hot-coffee/internal/service"
	"hot-coffee/models"
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
		help.WriteError(w, http.StatusInternalServerError, "Failed to fetch menu items")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *MenuHandler) AddNewMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.MenuService.AddNewMenuItem(item); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to add menu item")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		help.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.MenuService.UpdateMenuItem(item); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to update menu item")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	menuItemID := r.URL.Query().Get("menu_item_id")
	if menuItemID == "" {
		help.WriteError(w, http.StatusBadRequest, "Missing menu_item_id parameter")
		return
	}
	if err := h.MenuService.DeleteMenuItem(menuItemID); err != nil {
		help.WriteError(w, http.StatusInternalServerError, "Failed to delete menu item")
		return
	}
	w.WriteHeader(http.StatusOK)
}
