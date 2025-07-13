package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"net/http"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	total, err := h.service.GetTotalSales()
	if err != nil {
		http.Error(w, "Failed to get total sales", http.StatusInternalServerError)
		return
	}
	response := map[string]float64{"total_sales": total}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ReportHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetPopularItems()
	if err != nil {
		http.Error(w, "Failed to get popular items", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
