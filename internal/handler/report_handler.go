package handler

import (
	"encoding/json"
	"hot-coffee/help"
	"hot-coffee/internal/service"
	"log/slog"
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
		slog.Error("Failed to calculate total sales", "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to get total sales")
		return
	}
	slog.Info("Total sales calculated", "amount", total)
	response := map[string]float64{"total_sales": total}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ReportHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetPopularItems()
	if err != nil {
		slog.Error("Failed to generate popular items report", "error", err)
		help.WriteError(w, http.StatusInternalServerError, "Failed to get popular items")
		return
	}
	slog.Info("Popular items report generated", "count", len(items))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
