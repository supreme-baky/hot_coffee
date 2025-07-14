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
	dir := flag.String("dir", "data", "Path to the data directory")
	flag.Parse()

	if *helpFlag {
		help.PrintInfo()
		return
	}

	inventoryRepo := dal.NewJSONInventoryManager(filepath.Join(*dir, "inventory.json"))
	menuRepo := dal.NewJSONMenuManager(filepath.Join(*dir, "menu_items.json"))
	orderRepo := dal.NewJSONOrderManager(filepath.Join(*dir, "orders.json"))

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
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/inventory/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/inventory/")
		switch r.Method {
		case http.MethodGet:
			inventoryHandler.GetInventoryItem(w, r, id)
		case http.MethodPut:
			inventoryHandler.UpdateInventoryItem(w, r, id)
		case http.MethodDelete:
			inventoryHandler.DeleteInventoryItem(w, r, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	
	mux.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			menuHandler.GetAllMenuItems(w, r)
		case http.MethodPost:
			menuHandler.AddNewMenuItem(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/menu/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/menu/")
		switch r.Method {
		case http.MethodGet:
			menuHandler.GetMenuItem(w, r, id)
		case http.MethodPut:
			menuHandler.UpdateMenuItem(w, r, id)
		case http.MethodDelete:
			menuHandler.DeleteMenuItem(w, r, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderHandler.GetAllOrders(w, r)
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		if strings.HasSuffix(path, "/close") && r.Method == http.MethodPost {
			id := strings.TrimSuffix(path, "/close")
			orderHandler.CloseOrder(w, r, id)
			return
		}

		id := strings.TrimSuffix(path, "/")
		switch r.Method {
		case http.MethodGet:
			orderHandler.GetOrderByID(w, r, id)
		case http.MethodPut:
			orderHandler.UpdateOrder(w, r, id)
		case http.MethodDelete:
			orderHandler.DeleteOrder(w, r, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
