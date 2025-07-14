package main

import (
	"flag"
	"fmt"
	"hot-coffee/help"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	helpFlag := flag.Bool("help", false, "Prints help information")
	port := flag.Int("port", 8080, "Port number for the server")
	dirFlag := flag.String("dir", "./data", "Creates a directory with initial JSON files (orders, menu_items, inventory)")
	flag.Parse()

	if *helpFlag {
		help.PrintInfo()
		return
	}

	if *dirFlag != "" {
		if err := help.CreateDataDirWithFiles(*dirFlag); err != nil {
			log.Fatalf("Failed to create directory and files: %v", err)
		}
	}
	dataDir := *dirFlag

	inventoryRepo := dal.NewJSONInventoryManager(filepath.Join(dataDir, "inventory.json"))
	menuRepo := dal.NewJSONMenuManager(filepath.Join(dataDir, "menu_items.json"))
	orderRepo := dal.NewJSONOrderManager(filepath.Join(dataDir, "orders.json"))

	inventoryService := service.NewInventoryService(inventoryRepo)
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo, menuRepo, inventoryRepo)
	reportService := service.NewReportService(orderRepo, menuRepo)

	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	menuHandler := handler.NewMenuHandler(menuService)
	orderHandler := handler.NewOrderHandler(orderService)
	reportHandler := handler.NewReportHandler(reportService)

	mux := http.NewServeMux()

	mux.HandleFunc("/inventory", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			inventoryHandler.GetAllInventoryItems(w, r)
		case http.MethodPost:
			inventoryHandler.AddNewInventoryItem(w, r)
		default:
			help.WriteError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		}
	})
	mux.HandleFunc("/inventory/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/inventory/")
		if id == "" {
			help.WriteError(w, http.StatusBadRequest, "Missing inventory ID")
			return
		}
		switch r.Method {
		case http.MethodGet:
			inventoryHandler.GetInventoryItem(w, r, id)
		case http.MethodPut:
			inventoryHandler.UpdateInventoryItem(w, r, id)
		case http.MethodDelete:
			inventoryHandler.DeleteInventoryItem(w, r, id)
		default:
			help.WriteError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		}
	})

	mux.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			menuHandler.GetAllMenuItems(w, r)
		case http.MethodPost:
			menuHandler.AddNewMenuItem(w, r)
		default:
			help.WriteError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		}
	})
	mux.HandleFunc("/menu/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/menu/")
		if id == "" {
			help.WriteError(w, http.StatusBadRequest, "Missing menu item ID")
			return
		}
		switch r.Method {
		case http.MethodGet:
			menuHandler.GetMenuItem(w, r, id)
		case http.MethodPut:
			menuHandler.UpdateMenuItem(w, r, id)
		case http.MethodDelete:
			menuHandler.DeleteMenuItem(w, r, id)
		default:
			help.WriteError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		}
	})

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderHandler.GetAllOrders(w, r)
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		default:
			help.WriteError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		}
	})
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		if path == "" {
			help.WriteError(w, http.StatusNotFound, "Order path not found")
			return
		}

		if strings.HasSuffix(path, "/close") && r.Method == http.MethodPost {
			orderID := strings.TrimSuffix(path, "/close")
			orderID = strings.TrimSuffix(orderID, "/")
			orderHandler.CloseOrder(w, r, orderID)
			return
		}

		orderID := strings.TrimSuffix(path, "/")
		switch r.Method {
		case http.MethodGet:
			orderHandler.GetOrderByID(w, r, orderID)
		case http.MethodPut:
			orderHandler.UpdateOrder(w, r, orderID)
		case http.MethodDelete:
			orderHandler.DeleteOrder(w, r, orderID)
		default:
			help.WriteError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		}
	})

	mux.HandleFunc("/reports/total-sales", reportHandler.GetTotalSales)
	mux.HandleFunc("/reports/popular-items", reportHandler.GetPopularItems)

	if *port < 1 || *port > 65535 {
		log.Fatalf("Invalid port number: %d. Must be between 1 and 65535.", *port)
	}
	address := fmt.Sprintf(":%d", *port)
	log.Printf("Server started at http://localhost%s\n", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
